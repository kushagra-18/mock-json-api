package com.mock_json.api.helpers;

public class StringHelpers {

    /**
     * Removes all leading and trailing whitespace from a string, 
     * or the specified character.
     * @param str
     * @param ch
     * @return
     */
    public static String ltrim(String str, char ch) {
        StringBuilder result = new StringBuilder(str);
        while (result.length() > 0 && result.charAt(0) == ch) {
            result.deleteCharAt(0);
        }
        return result.toString();
    }
}