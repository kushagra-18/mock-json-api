package controllers

import (
	// "bytes" // Removed unused import
	"database/sql" // Added for requestLog.UrlID
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v4" // Not directly used in controller if middleware handles it
	"gorm.io/gorm"

	"mockapi/config"
	"mockapi/dtos"
	"mockapi/models"
	"mockapi/services"
	"mockapi/utils"
)

// MockContentController handles API endpoints related to creating and serving mock content.
type MockContentController struct {
	projectService     *services.ProjectService
	mockContentService *services.MockContentService
	urlService         *services.URLService
	requestLogService  *services.RequestLogService
	redisService       *services.RedisService
	proxyService       *services.ProxyService // Added proxyService
	jwtSecret          string
	config             config.Config
}

// NewMockContentController creates a new MockContentController.
func NewMockContentController(
	projService *services.ProjectService,
	mcService *services.MockContentService,
	uService *services.URLService,
	rlService *services.RequestLogService,
	rService *services.RedisService,
	pService *services.ProxyService, // Added proxyService
	cfg config.Config,
) *MockContentController {
	return &MockContentController{
		projectService:     projService,
		mockContentService: mcService,
		urlService:         uService,
		requestLogService:  rlService,
		redisService:       rService,
		proxyService:       pService, // Added proxyService
		jwtSecret:          cfg.JWTSecretKey, // Store JWT secret from config
		config:             cfg,
	}
}

// SaveMockContent handles POST /mock/:projectSlug
func (mcc *MockContentController) SaveMockContent(c *gin.Context) {
	projectSlug := c.Param("projectSlug")

	var dto dtos.MockContentUrlDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	project, err := mcc.projectService.GetProjectBySlug(projectSlug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Project with slug '%s' not found.", projectSlug))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error finding project: "+err.Error())
		}
		return
	}

	existingURL, err := mcc.urlService.FindByProjectIDAndURL(project.ID, dto.URLData.URL)
	if err != nil && err != gorm.ErrRecordNotFound {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error checking for existing URL: "+err.Error())
		return
	}
	if existingURL != nil {
		utils.ErrorResponse(c, http.StatusConflict, fmt.Sprintf("URL path '%s' already exists for this project.", dto.URLData.URL))
		return
	}

	newURL := &models.Url{
		ProjectID:   project.ID,
		Name:        dto.URLData.Name,
		Description: utils.StringPointerToString(dto.URLData.Description),
		URL:         dto.URLData.URL,
		Status:      dto.URLData.Status,
	}

	if err := mcc.urlService.CreateURL(newURL, project.ID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create URL: "+err.Error())
		return
	}

	var mockContentsToSave []models.MockContent
	for _, mcDto := range dto.MockContentList {
		content := models.MockContent{
			UrlID:       newURL.ID,
			Name:        mcDto.Name,
			Description: utils.StringPointerToString(mcDto.Description),
			Data:        mcDto.Data,
			Randomness:  utils.Int64PointerToInt64(mcDto.Randomness),
			Latency:     utils.Int64PointerToInt64(mcDto.Latency),
		}
		mockContentsToSave = append(mockContentsToSave, content)
	}

	if len(mockContentsToSave) > 0 {
		savedMCs, err := mcc.mockContentService.SaveMockContentList(mockContentsToSave, newURL.ID)
		if err != nil {
			_ = mcc.urlService.DeleteURL(newURL.ID) // Attempt to clean up, ignore error for now
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save mock contents: "+err.Error())
			return
		}
		newURL.MockContents = savedMCs
	}

	responseDTO := struct {
		URL models.Url `json:"url"`
	} {
		URL: *newURL,
	}
	utils.SuccessResponse(c, http.StatusCreated, responseDTO)
}

// UpdateMockContent handles PATCH /mock/:projectSlug/:urlId
func (mcc *MockContentController) UpdateMockContent(c *gin.Context) {
	projectSlug := c.Param("projectSlug")
	urlIDStr := c.Param("urlId")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid URL ID format.")
		return
	}

	var dto dtos.UpdateMockContentUrlDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	project, err := mcc.projectService.GetProjectBySlug(projectSlug)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Project with slug '%s' not found.", projectSlug))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error finding project: "+err.Error())
		}
		return
	}

	urlToUpdate, err := mcc.urlService.GetURLByID(uint(urlID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("URL with ID %d not found.", urlID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error finding URL: "+err.Error())
		}
		return
	}
	if urlToUpdate.ProjectID != project.ID {
		utils.ErrorResponse(c, http.StatusForbidden, "URL does not belong to the specified project.")
		return
	}

	var mockContentsToUpdate []models.MockContent
	for _, mcDto := range dto.MockContentList {
		content := models.MockContent{
			UrlID:       uint(urlID),
			Name:        utils.StringPointerToString(mcDto.Name),
			Description: utils.StringPointerToString(mcDto.Description),
			Data:        utils.StringPointerToString(mcDto.Data),
			Randomness:  utils.Int64PointerToInt64(mcDto.Randomness),
			Latency:     utils.Int64PointerToInt64(mcDto.Latency),
		}
		if mcDto.ID != nil {
			content.ID = *mcDto.ID
			content.BaseModel.ID = *mcDto.ID
		}
		mockContentsToUpdate = append(mockContentsToUpdate, content)
	}

	updatedMCs, err := mcc.mockContentService.UpdateMockContentList(mockContentsToUpdate, uint(urlID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update mock contents: "+err.Error())
		return
	}

	urlToUpdate.MockContents = updatedMCs
	utils.SuccessResponse(c, http.StatusOK, urlToUpdate)
}

// GetMockedJSON handles GET /mock/:teamSlug/:projectSlug/*wildcardPath
func (mcc *MockContentController) GetMockedJSON(c *gin.Context) {
	teamSlug := c.Param("teamSlug")
	projectSlug := c.Param("projectSlug")
	encodedPathParams := c.Param("wildcardPath")

	requestLog := &models.RequestLog{
		IPAddress: c.ClientIP(),
		Timestamp: time.Now(),
		Method:    c.Request.Method,
		URL:       c.Request.URL.String(),
		IsProxied: false,
		CreatedAt: time.Now(),
	}

	pathParts := strings.SplitN(strings.TrimPrefix(encodedPathParams, "/"), "/", 2)
	if len(pathParts) < 1 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid path structure after project slug.")
		mcc.finalizeRequestLog(requestLog, http.StatusBadRequest, 0, 0)
		return
	}

	actualPath := ""
	var decodedParams dtos.GetMockedJSONParamsDTO
	potentialBase64 := pathParts[0]
	decodedBytes, err := base64.URLEncoding.DecodeString(potentialBase64)

	if err == nil {
		if errJson := json.Unmarshal(decodedBytes, &decodedParams); errJson != nil {
			log.Printf("WARN: Failed to unmarshal decoded base64 params: %v. Raw: %s", errJson, string(decodedBytes))
		}
		if len(pathParts) > 1 {
			actualPath = "/" + pathParts[1]
		} else {
			actualPath = "/"
		}
	} else {
		actualPath = encodedPathParams
		if !strings.HasPrefix(actualPath, "/") {
			actualPath = "/" + actualPath
		}
	}
	requestLog.URL = actualPath

	globalRateLimitKey := mcc.redisService.CreateRedisKey("ratelimit:global", c.ClientIP())
	isGloballyLimited, rlErr := mcc.redisService.RateLimit(globalRateLimitKey, mcc.config.GlobalMaxAllowedRequests, int64(mcc.config.GlobalTimeWindowSeconds))
	if rlErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Error checking global rate limit.")
		mcc.finalizeRequestLog(requestLog, http.StatusInternalServerError, 0, 0)
		return
	}
	if isGloballyLimited {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "Global rate limit exceeded.")
		mcc.finalizeRequestLog(requestLog, http.StatusTooManyRequests, 0, 0)
		return
	}

	project, err := mcc.projectService.GetProjectByTeamSlugAndProjectSlug(teamSlug, projectSlug)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound { statusCode = http.StatusNotFound }
		utils.ErrorResponse(c, statusCode, "Project not found or error fetching project: "+err.Error())
		mcc.finalizeRequestLog(requestLog, statusCode, 0, 0)
		return
	}
	requestLog.ProjectID = project.ID

	isForwardCall := decodedParams.IsForwardCall != nil && *decodedParams.IsForwardCall
	userWantsForward := decodedParams.Forward != nil && *decodedParams.Forward

	if project.IsForwardProxyActive && userWantsForward && !isForwardCall {
		proxySettings, err := mcc.proxyService.GetForwardProxyByProjectID(project.ID)
		if err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error fetching proxy settings.")
			mcc.finalizeRequestLog(requestLog, http.StatusInternalServerError, project.ID, 0)
			return
		}
		if proxySettings != nil && proxySettings.Domain != "" {
			requestLog.IsProxied = true
			targetURL := mcc.buildTargetURL(proxySettings.Domain, c.Request)

			newDecodedParams := decodedParams
			newDecodedParams.IsForwardCall = utils.BoolPointer(true)
			newParamsBytes, _ := json.Marshal(newDecodedParams)
			newBase64Params := base64.URLEncoding.EncodeToString(newParamsBytes)
			proxiedPath := fmt.Sprintf("/mock/%s/%s/%s%s", teamSlug, projectSlug, newBase64Params, actualPath)

			httpClient := &http.Client{Timeout: 10 * time.Second}
			req, err := http.NewRequest(c.Request.Method, targetURL+proxiedPath, c.Request.Body)
			if err != nil {
				utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create proxy request: "+err.Error())
				mcc.finalizeRequestLog(requestLog, http.StatusInternalServerError, project.ID, 0)
				return
			}
			req.Header = c.Request.Header.Clone()

			resp, err := httpClient.Do(req)
			if err != nil {
				utils.ErrorResponse(c, http.StatusBadGateway, "Failed to execute proxy request: "+err.Error())
				mcc.finalizeRequestLog(requestLog, http.StatusBadGateway, project.ID, 0)
				return
			}
			defer resp.Body.Close()

			for key, values := range resp.Header {
				for _, value := range values {
					c.Writer.Header().Add(key, value)
				}
			}
			c.Writer.WriteHeader(resp.StatusCode)
			_, _ = io.Copy(c.Writer, resp.Body)
			mcc.finalizeRequestLog(requestLog, resp.StatusCode, project.ID, 0)
			return
		}
	}

	urlData, err := mcc.urlService.GetURLByTeamSlugProjectSlugAndPath(teamSlug, projectSlug, actualPath)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == gorm.ErrRecordNotFound { statusCode = http.StatusNotFound }
		utils.ErrorResponse(c, statusCode, "URL not found or error fetching URL: "+err.Error())
		mcc.finalizeRequestLog(requestLog, statusCode, project.ID, 0)
		return
	}
	requestLog.UrlID = sql.NullInt64{Int64: int64(urlData.ID), Valid: true}

	if len(urlData.MockContents) == 0 {
		utils.ErrorResponse(c, http.StatusNotFound, "No mock content available for this URL.")
		mcc.finalizeRequestLog(requestLog, http.StatusNotFound, project.ID, urlData.ID)
		return
	}
	selectedMock := mcc.mockContentService.SelectRandomMockContent(urlData.MockContents)
	if selectedMock == nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to select mock content.")
		mcc.finalizeRequestLog(requestLog, http.StatusInternalServerError, project.ID, urlData.ID)
		return
	}

	mcc.mockContentService.SimulateLatency(selectedMock.Latency)
	_ = mcc.urlService.IncrementRequestStats(urlData.ID)

	var jsonOutput interface{}
	responseStatusCode := mcc.getStatusCodeInt(urlData.Status) // Use helper for status code

	if err := json.Unmarshal([]byte(selectedMock.Data), &jsonOutput); err != nil {
		c.Data(responseStatusCode, "text/plain; charset=utf-8", []byte(selectedMock.Data)) // Corrected charset
	} else {
		c.JSON(responseStatusCode, jsonOutput)
	}
	mcc.finalizeRequestLog(requestLog, responseStatusCode, project.ID, urlData.ID)
}

func (mcc *MockContentController) finalizeRequestLog(logEntry *models.RequestLog, statusCode int, projectID uint, urlID uint) {
	logEntry.Status = statusCode
	if projectID != 0 {
		logEntry.ProjectID = projectID
	}
	if urlID != 0 {
		logEntry.UrlID = sql.NullInt64{Int64: int64(urlID), Valid: true}
	}

	if err := mcc.requestLogService.SaveRequestLog(logEntry); err != nil {
		log.Printf("ERROR: Failed to save request log: %v. Entry: %+v", err, logEntry)
	}
}

func (mcc *MockContentController) buildTargetURL(baseDomain string, originalReq *http.Request) string {
	scheme := "http"
	if originalReq.TLS != nil || strings.EqualFold(originalReq.Header.Get("X-Forwarded-Proto"), "https") {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, baseDomain)
}

// Helper for DTO optional fields
func StringPointerToString(s *string) string {
	if s == nil { return "" }
	return *s
}

func Int64PointerToInt64(i *int64) int64 {
	if i == nil { return 0 }
	return *i
}
func BoolPointer(b bool) *bool { return &b }


// getStatusCodeInt converts models.StatusCode to an int.
// This should ideally be a method on models.StatusCode or defined in the models package.
func (mcc *MockContentController) getStatusCodeInt(statusCode models.StatusCode) int {
	// This is a simplified map. A more robust solution would be in models/common.go
	// and cover all defined status codes.
	statusMapping := map[models.StatusCode]int{
		models.StatusOK:                   http.StatusOK,
		models.StatusCreated:              http.StatusCreated,
		models.StatusAccepted:             http.StatusAccepted,
		models.StatusNonAuthoritativeInfo: 203, // http.StatusNonAuthoritativeInformation,
		models.StatusNoContent:            http.StatusNoContent,
		models.StatusResetContent:         http.StatusResetContent,
		models.StatusPartialContent:       http.StatusPartialContent,
		models.StatusMovedPermanently:     http.StatusMovedPermanently,
		models.StatusFound:                http.StatusFound,
		models.StatusSeeOther:             http.StatusSeeOther,
		models.StatusNotModified:          http.StatusNotModified,
		models.StatusTemporaryRedirect:    http.StatusTemporaryRedirect,
		models.StatusPermanentRedirect:    http.StatusPermanentRedirect,
		models.StatusBadRequest:           http.StatusBadRequest,
		models.StatusUnauthorized:         http.StatusUnauthorized,
		models.StatusPaymentRequired:      http.StatusPaymentRequired,
		models.StatusForbidden:            http.StatusForbidden,
		models.StatusNotFound:             http.StatusNotFound,
		models.StatusMethodNotAllowed:     http.StatusMethodNotAllowed,
		models.StatusNotAcceptable:        http.StatusNotAcceptable,
		models.StatusProxyAuthRequired:    http.StatusProxyAuthRequired,
		models.StatusRequestTimeout:       http.StatusRequestTimeout,
		models.StatusConflict:             http.StatusConflict,
		models.StatusGone:                 http.StatusGone,
		models.StatusLengthRequired:       http.StatusLengthRequired,
		models.StatusPreconditionFailed:   http.StatusPreconditionFailed,
		models.StatusRequestEntityTooLarge: http.StatusRequestEntityTooLarge,
		models.StatusRequestURITooLong:    http.StatusRequestURITooLong,
		models.StatusUnsupportedMediaType: http.StatusUnsupportedMediaType,
		models.StatusRequestedRangeNotSatisfiable: http.StatusRequestedRangeNotSatisfiable,
		models.StatusExpectationFailed:    http.StatusExpectationFailed,
		models.StatusTeapot:               http.StatusTeapot,
		models.StatusUnprocessableEntity:  http.StatusUnprocessableEntity,
		models.StatusLocked:               http.StatusLocked,
		models.StatusFailedDependency:     http.StatusFailedDependency,
		models.StatusTooEarly:             http.StatusTooEarly,
		models.StatusUpgradeRequired:      http.StatusUpgradeRequired,
		models.StatusPreconditionRequired: http.StatusPreconditionRequired,
		models.StatusTooManyRequests:      http.StatusTooManyRequests,
		models.StatusRequestHeaderFieldsTooLarge: http.StatusRequestHeaderFieldsTooLarge,
		models.StatusUnavailableForLegalReasons: http.StatusUnavailableForLegalReasons,
		models.StatusInternalServerError:  http.StatusInternalServerError,
		models.StatusNotImplemented:       http.StatusNotImplemented,
		models.StatusBadGateway:           http.StatusBadGateway,
		models.StatusServiceUnavailable:   http.StatusServiceUnavailable,
		models.StatusGatewayTimeout:       http.StatusGatewayTimeout,
		models.StatusHTTPVersionNotSupported: http.StatusHTTPVersionNotSupported,
		models.StatusVariantAlsoNegotiates: http.StatusVariantAlsoNegotiates,
		models.StatusInsufficientStorage:  http.StatusInsufficientStorage,
		models.StatusLoopDetected:         http.StatusLoopDetected,
		models.StatusNotExtended:          http.StatusNotExtended,
		models.StatusNetworkAuthenticationRequired: http.StatusNetworkAuthenticationRequired,
	}
	val, ok := statusMapping[statusCode]
	if !ok {
		log.Printf("Warning: Unmapped StatusCode '%s', defaulting to 200 OK", statusCode)
		return http.StatusOK // Default if not found
	}
	return val
}

// Note: The local `statusMap` and `(sc models.StatusCode) Code()` method previously in this file
// were removed in favor of `getStatusCodeInt` for clarity and to consolidate the mapping logic.
// The `sql.NullInt64` import was added for `requestLog.UrlID`.
// `github.com/golang-jwt/jwt/v4` import was commented out as it's not directly used if middleware handles JWT.
// Charset for text/plain response in GetMockedJSON was corrected to "utf-8".
// Error handling for `mcc.urlService.DeleteURL` and `mcc.urlService.IncrementRequestStats` was made less verbose (ignoring error for cleanup).
// `io.Copy` result was ignored in proxy logic.
