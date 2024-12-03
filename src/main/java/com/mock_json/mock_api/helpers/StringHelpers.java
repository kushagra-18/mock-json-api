package com.mock_json.mock_api.helpers;

import java.util.concurrent.ThreadLocalRandom;

public class StringHelpers {

    /**
     * Removes all leading and trailing whitespace from a string,
     * or the specified character.
     * 
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

    /**
     * Generates a random string
     * 
     * @param length
     * @return
     */
    public static String generateRandomString(int length) {
        StringBuilder result = new StringBuilder(length);
        for (int i = 0; i < length; i++) {
            // Generate a random number between 'a' (97) and 'z' (122)
            char randomChar = (char) ThreadLocalRandom.current().nextInt('a', 'z' + 1);
            result.append(randomChar);
        }
        return result.toString();
    }

    /**
     * Converts a slug to a human-readable string.
     * @param slug the slug to be converted (e.g., "hello-world")
     * @return the human-readable string (e.g., "Hello World")
     */
    
    public static String unslug(String slug) {
        if (slug == null || slug.isEmpty()) {
            return "";
        }

        // Split the slug by dashes or underscores
        String[] words = slug.split("[-_]+");

        // Capitalize the first letter of each word and join them with a space
        StringBuilder readable = new StringBuilder();
        for (String word : words) {
            if (!word.isEmpty()) {
                readable.append(Character.toUpperCase(word.charAt(0)))
                        .append(word.substring(1).toLowerCase())
                        .append(" ");
            }
        }

        // Trim the trailing space and return
        return readable.toString().trim();
    }
}