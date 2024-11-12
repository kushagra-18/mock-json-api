package com.mock_json.mock_api.contexts;

public class HeaderContext {

    private static ThreadLocal<String> teamSlug = new ThreadLocal<>();
    private static ThreadLocal<String> projectSlug = new ThreadLocal<>();

    public static void setTeamSlug(String team) {
        teamSlug.set(team);
    }

    public static String getTeamSlug() {
        return teamSlug.get();
    }

    public static void setProjectSlug(String project) {
        projectSlug.set(project);
    }

    public static String getProjectSlug() {
        return projectSlug.get();
    }

    public static void clear() {
        teamSlug.remove();
        projectSlug.remove();
    }
}
