// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package com.example.bnb.catalog;

import org.joda.time.DateTime;
import org.joda.time.format.ISODateTimeFormat;
import ch.qos.logback.core.encoder.EncoderBase;
import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.spi.ILoggingEvent;
import java.util.HashMap;
import java.util.Map;

import com.google.gson.Gson;

public class LoggingEventGoogleCloudEncoder extends EncoderBase<ILoggingEvent> {
    private static final byte[] EMPTY_BYTES = new byte[0];
    private final Gson gson = new Gson();

    @Override
    public byte[] headerBytes() {
        return EMPTY_BYTES;
    }

    @Override
    public byte[] encode(ILoggingEvent e) {
        DateTime timestamp = new DateTime(e.getTimeStamp());
        Map<String, Object> fields = new HashMap<String, Object>() {
            {
                put("timestamp", ISODateTimeFormat.dateTime().print(timestamp));
                put("severity", severityFor(e.getLevel()));
                put("message", e.getMessage());
            }
        };
        String data = gson.toJson(fields) + "\n";
        return data.getBytes();
    }

    @Override
    public byte[] footerBytes() {
        return EMPTY_BYTES;
    }

    private static String severityFor(Level level) {
        switch (level.toInt()) {
            case Level.TRACE_INT:
                return "DEBUG";
            case Level.DEBUG_INT:
                return "DEBUG";
            case Level.INFO_INT:
                return "INFO";
            case Level.WARN_INT:
                return "WARNING";
            case Level.ERROR_INT:
                return "ERROR";
            default:
                return "DEFAULT";
        }
    }
}