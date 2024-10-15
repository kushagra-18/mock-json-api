package com.mock_json.api.contexts;

public class HeaderContext {

    private static ThreadLocal<String> teamId = new ThreadLocal<>();
    private static ThreadLocal<String> projectId = new ThreadLocal<>();

    public static void setTeamId(String team) {
        teamId.set(team);
    }

    public static String getTeamId() {
        return teamId.get();
    }

    public static void setProjectId(String project) {
        projectId.set(project);
    }

    public static String getProjectId() {
        return projectId.get();
    }

    public static void clear() {
        teamId.remove();
        projectId.remove();
    }
}
