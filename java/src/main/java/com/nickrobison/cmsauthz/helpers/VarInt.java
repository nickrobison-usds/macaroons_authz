package com.nickrobison.cmsauthz.helpers;

import java.nio.ByteBuffer;

/*
 * Copyright 2016 Kukri Máté.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

public class VarInt {
    /* bitmasks for splitting int */
    private static final int[] bitmasks = { 0x0000007F, 0x00003F80, 0x001FC000, 0x0FE00000, 0xF0000000 };

    /* shift value for nibbles */
    private static final int[] shift_value = { 0, 7, 14, 21, 28 };

    public static ByteBuffer encodeVarInt(int input) {
        /* do signed magic */
        int value = (input << 1) ^ (input >> 31);
        ByteBuffer buf = ByteBuffer.allocate(5);

        for (int i = 0; i < 5; ++i) {
            byte b = (byte) ((bitmasks[i] & value) >> shift_value[i]);

            byte nextb = (byte) ((bitmasks[i+1] & value) >> shift_value[i+1]);
            if (nextb == 0) {
                buf.put(b);
                break;
            }

            b = (byte) (0x80 | b);
            buf.put(b);
        }

        buf.flip();
        return buf;
    }

    public static int decodeVarInt(ByteBuffer varint) {
        int out = 0;
        for (int i = 0; i < varint.limit(); ++i) {
            byte b = varint.get(i);
            if (i + 1 != varint.limit()) {
                b = (byte) (0x80 ^ b);
            }

            out |= b << shift_value[i];
        }

        /* undo signed magic */
        return (out >>> 1) ^ -(out & 1);
    }
}