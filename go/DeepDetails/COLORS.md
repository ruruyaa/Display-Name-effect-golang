# 🎨 Known Colors Index — Golang

Discord Display Name Styles require color representations in **24-bit decimal integer format**, rather than standard hex strings (e.g. `"#FF00FF"`). This document maps color options to their decimal values and details how translation functions work in Go.

---

## 📋 Standard Decimal Colors List

Here is your core reference for pre-registered hues:

| Color Name | Hex Representation | Decimal Integer | Example Use-Case |
|---|---|---|---|
| **White** | `#FFFFFF` | `16777215` | Default elegant overlay accent |
| **Blue** | `#0016E9` | `5865` | High-tech deep blue accent |
| **Pink** | `#FF00FF` | `16711935` | Vibrant aesthetic neon glow |
| **Purple** | `#800000` | `8388736` | Heavy dark gothic contrast |

---

## 📈 Multi-Color / Gradient Pair Presets

When combining values for gradients (`effect_id: 2`) or multi-colored glow overlays, supply the integer values in order inside the payload slice:

### 1. White to Blue Gradient
- **Hex Code Pair**: `#FFFFFF` + `#0016E9`
- **Decimal Representation**: `[]int{16777215, 5865}`
- **Style Concept**: Cold, electronic glass effect. Pair with font ID `7` (Neo-Castel).

### 2. Pink to Purple Gradient
- **Hex Code Pair**: `#FF00FF` + `#800000`
- **Decimal Representation**: `[]int{16711935, 8388736}`
- **Style Concept**: Cyberpunk high-saturation contrast. Pair with font ID `1` (Bangers) or `9` (Ribes).

### 3. Pure Cyberpunk Glow
- **Hex Code Pair**: `#00FFD2` (Cyan) + `#FF00FF` (Pink)
- **Decimal Representation**: `[]int{65490, 16711935}`
- **Style Concept**: Distinctive neon glow aesthetic. Pair with font ID `8` (Pixelify Sans) or `10` (Sinistre).

---

## 💻 Golang Conversion Functions

You can integrate these helpers directly inside your Go packages to safely convert on-the-fly. Notice how Go returns `error` types instead of throwing exceptions to keep your server completely crash-proof.

```go
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// HexToDecimal converts a standard Hex string (e.g., "#FF00FF" or "FF00FF")
// to a Discord-compatible 24-bit decimal integer.
func HexToDecimal(hexCode string) (int, error) {
	sanitized := strings.TrimSpace(strings.ReplaceAll(hexCode, "#", ""))
	
	parsed, err := strconv.ParseInt(sanitized, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hexadecimal color input: %s", hexCode)
	}

	if parsed < 0 || parsed > 0xffffff {
		return 0, fmt.Errorf("hex color code out of 24-bit bounds: %s", hexCode)
	}
	
	return int(parsed), nil
}

// DecimalToHex converts a Discord-compatible 24-bit decimal integer
// back to a formatted Hex string.
func DecimalToHex(decimal int) (string, error) {
	if decimal < 0 || decimal > 0xffffff {
		return "", fmt.Errorf("invalid decimal color input: %d", decimal)
	}

	// The %06X verb formats the integer as an uppercase hex string, 
	// padded with leading zeros to ensure it is always 6 characters long.
	return fmt.Sprintf("#%06X", decimal), nil
}

/* // Example usage:
func main() {
	decColor, err := HexToDecimal("#FF00FF") 
	if err == nil {
		fmt.Println(decColor) // Prints: 16711935
	}

	hexColor, err := DecimalToHex(5865)
	if err == nil {
		fmt.Println(hexColor) // Prints: #0016E9
	}
}
*/

