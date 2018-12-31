package com.nickrobison.cmsauthz.helpers;

import java.util.Arrays;

public class Helpers {

    private Helpers() {
//        Not used
    }

    /**
     * Simple helper method for printing bytes to the console, while converting them to unsigned.
     * This makes it a bit easier to debug encodings between languages which don't use signed byte values.
     *
     * @param name        - {@link String} to display before printing the bytes values. Useful to help differentiate the outputs.
     * @param signedBytes - {@link byte[]} to convert to unsigned and display
     */
    public static void printUnsignedBytes(String name, byte[] signedBytes) {

        int[] unsigned = new int[signedBytes.length];

        for (int i = 0; i < signedBytes.length; i++) {
            unsigned[i] = (signedBytes[i] & 0xFF);
        }

        System.out.printf("%s as unsigned bytes:\n", name);
        System.out.println(Arrays.toString(unsigned));
    }
}
