#!/bin/bash -eu
# Copyright 2025 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SCRIPT_DIR="$(realpath "$(dirname "${0}")")"; readonly SCRIPT_DIR
BUCKET_ID="$1"
if [ -z "${BUCKET_ID}" ]; then
    echo "Please provide valid bucket id in the form 'gs://bucket_id'"
    exit 1
fi

pushd "${SCRIPT_DIR}../images"

gcloud storage cp mug.jpg loafers.jpg hairdryer.jpg "${BUCKET_ID}/listing/listing1/"
gcloud storage cp watch.jpg candle-holder.jpg "${BUCKET_ID}/listing/listing2/"
gcloud storage cp sunglasses.jpg tank-top.jpg bamboo-glass-jar.jpg "${BUCKET_ID}/listing/listing3/"

popd