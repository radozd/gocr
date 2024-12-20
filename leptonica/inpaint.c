/*
PIX *pixInpaint(PIX *pixs, PIX *mask, l_int32 radius) {
    l_int32 width = pixGetWidth(pixs);
    l_int32 height = pixGetHeight(pixs);
    PIX *pixd = pixCreate(width, height, 32); // Create a new image for the output
    l_uint32 *data = (l_uint32 *)pixGetData(pixs);
    l_uint32 *maskData = (l_uint32 *)pixGetData(mask);
    l_uint32 *outputData = (l_uint32 *)pixGetData(pixd);

    // Iterate over each pixel in the source image
    for (l_int32 y = 0; y < height; y++) {
        for (l_int32 x = 0; x < width; x++) {
            // Check if the current pixel is part of the mask
            if (GET_DATA_BIT(maskData, x, y)) {
                // Initialize variables to accumulate color values
                l_int32 r = 0, g = 0, b = 0, count = 0;

                // Check neighboring pixels within the specified radius
                for (l_int32 dy = -radius; dy <= radius; dy++) {
                    for (l_int32 dx = -radius; dx <= radius; dx++) {
                        l_int32 nx = x + dx;
                        l_int32 ny = y + dy;

                        // Ensure the neighbor is within bounds
                        if (nx >= 0 && nx < width && ny >= 0 && ny < height) {
                            // Check if the neighbor is not part of the mask
                            if (!GET_DATA_BIT(maskData, nx, ny)) {
                                l_uint32 pixel = data[ny * width + nx];
                                r += (pixel >> 16) & 0xFF; // Red
                                g += (pixel >> 8) & 0xFF;  // Green
                                b += pixel & 0xFF;         // Blue
                                count++;
                            }
                        }
                    }
                }

                // If we found valid neighbors, calculate the average color
                if (count > 0) {
                    r /= count;
                    g /= count;
                    b /= count;
                    outputData[y * width + x] = (r << 16) | (g << 8) | b; // Set the pixel color
                } else {
                    outputData[y * width + x] = 0; // Default to black if no neighbors found
                }
            } else {
                // If not part of the mask, copy the original pixel
                outputData[y * width + x] = data[y * width + x];
            }
        }
    }

    return pixd; // Return the inpainted image
}

PIX *pixInpaintTelea(PIX *pixs, PIX *mask, l_int32 radius) {
    l_int32 width = pixGetWidth(pixs);
    l_int32 height = pixGetHeight(pixs);
    PIX *pixd = pixCreate(width, height, 32); // Create a new image for the output
    l_uint32 *data = (l_uint32 *)pixGetData(pixs);
    l_uint32 *maskData = (l_uint32 *)pixGetData(mask);
    l_uint32 *outputData = (l_uint32 *)pixGetData(pixd);

    // Initialize the output image with the original image
    for (l_int32 y = 0; y < height; y++) {
        for (l_int32 x = 0; x < width; x++) {
            outputData[y * width + x] = data[y * width + x];
        }
    }

    // Iterate over each pixel in the mask
    for (l_int32 y = 0; y < height; y++) {
        for (l_int32 x = 0; x < width; x++) {
            // Check if the current pixel is part of the mask
            if (GET_DATA_BIT(maskData, x, y)) {
                // Variables to accumulate color values and weights
                double r = 0, g = 0, b = 0, totalWeight = 0;

                // Check neighboring pixels within the specified radius
                for (l_int32 dy = -radius; dy <= radius; dy++) {
                    for (l_int32 dx = -radius; dx <= radius; dx++) {
                        l_int32 nx = x + dx;
                        l_int32 ny = y + dy;

                        // Ensure the neighbor is within bounds
                        if (nx >= 0 && nx < width && ny >= 0 && ny < height) {
                            // Check if the neighbor is not part of the mask
                            if (!GET_DATA_BIT(maskData, nx, ny)) {
                                l_uint32 pixel = data[ny * width + nx];
                                double weight = 1.0 / (1.0 + sqrt(dx * dx + dy * dy)); // Weight based on distance

                                r += ((pixel >> 16) & 0xFF) * weight; // Red
                                g += ((pixel >> 8) & 0xFF) * weight;  // Green
                                b += (pixel & 0xFF) * weight;         // Blue
                                totalWeight += weight;
                            }
                        }
                    }
                }

                // If we found valid neighbors, calculate the average color
                if (totalWeight > 0) {
                    outputData[y * width + x] = ((l_uint32)(r / totalWeight) << 16) |
                                                  ((l_uint32)(g / totalWeight) << 8) |
                                                  (l_uint32)(b / totalWeight); // Set the pixel color
                }
            }
        }
    }

    return pixd; // Return the inpainted image
}

*/