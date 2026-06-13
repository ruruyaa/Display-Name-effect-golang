# 🔤 Known Fonts Reference Guide — Golang

This file catalogs every known typography style currently configured inside the **Discord Display Name Styles** service. It serves as a definitive reference mapping for styling presets using Go structs.

---

## 📋 Comprehensive Fonts Catalog

Each font listed below carries a specific integer ID used inside payload integrations (Format A: `font_id`, Format B: `display_name_font_id`).

### 1. Bangers
- **ID**: `1`
- **Key Name**: `bangers`
- **Description**: Bold, wide, comic-style display font with tall caps. Perfect for a punchy, active, or playful aesthetic.
- **Example Usage**: `style := Style{FontID: 1, EffectID: 1, Colors: []int{16711935}}`
- **Compatibility Status**: Supported on all platforms. Best paired with the "Pop" or "Neon" effects.

### 2. BioRhyme
- **ID**: `2`
- **Key Name**: `biorhyme`
- **Description**: Elegant, expansive, heavy slab-serif font utilizing modern wide-proportional tracking.
- **Example Usage**: `style := Style{FontID: 2, EffectID: 1, Colors: []int{16777215}}`
- **Compatibility Status**: Experimental. Some mobile clients render standard sans-serif fallback if scaling is limited.

### 3. Cherry Bomb (Sakura)
- **ID**: `3`
- **Key Name**: `cherry_bomb`
- **Description**: Highly playful, rounded bubble-letter typography. Excellent for cute, casual, or community-centric bot designs.
- **Example Usage**: `style := Style{FontID: 3, EffectID: 4, Colors: []int{16777215}}` (Often paired with Toon style).
- **Compatibility Status**: High. Popular on modern server profile rollouts.

### 4. Chicle (Jellybean)
- **ID**: `4`
- **Key Name**: `chicle`
- **Description**: Rounded, soft, curved, jellybean-style display letters.
- **Example Usage**: `style := Style{FontID: 4, EffectID: 6, Colors: []int{16777215}}`
- **Compatibility Status**: High. Safe for all profile layouts.

### 5. Compagnon
- **ID**: `5`
- **Key Name**: `compagnon`
- **Description**: Historic, monospaced typewriter font. Conveys a rustic, scientific, or highly technical aesthetic.
- **Example Usage**: `style := Style{FontID: 5, EffectID: 1, Colors: []int{16777215}}`
- **Compatibility Status**: Experimental. Rendering might feel small when squeezed inside user lists on mobile devices.

### 6. Museo Moderno
- **ID**: `6`
- **Key Name**: `museo_moderno`
- **Description**: Ultra-modern, rounded, circular geometric sans-serif typeface. Looks stellar for clean, high-tech, or futuristic layouts.
- **Example Usage**: `style := Style{FontID: 6, EffectID: 3, Colors: []int{5865}}`
- **Compatibility Status**: High.

### 7. Neo-Castel (Medieval)
- **ID**: `7`
- **Key Name**: `neo_castel`
- **Description**: Gothic, medieval-style display font mimicking calligraphy. Imparts a fantasy, ancient, or dramatic aura.
- **Example Usage**: `style := Style{FontID: 7, EffectID: 2, Colors: []int{5865, 16777215}}` (Stellar with color gradients).
- **Compatibility Status**: Moderate.

### 8. Pixelify Sans (8Bit)
- **ID**: `8`
- **Key Name**: `pixelify_sans`
- **Description**: retro, blocky, 8-bit arcade style pixel-art font. Ideal for gamer-centric guilds or nostalgic designs.
- **Example Usage**: `style := Style{FontID: 8, EffectID: 5, Colors: []int{8388736}}`
- **Compatibility Status**: Extremely High. Fully verified on desktop and mobile clients.

### 9. Ribes
- **ID**: `9`
- **Key Name**: `ribes`
- **Description**: Stylish, italicized geometric brush typeface. Bold display structure with high slant curves.
- **Example Usage**: `style := Style{FontID: 9, EffectID: 3, Colors: []int{16711935}}` (Pink neon).
- **Compatibility Status**: Verified. High usage across elite presets.

### 10. Sinistre (Vampyre)
- **ID**: `10`
- **Key Name**: `sinistre`
- **Description**: Dark, sharp, vampyric, elegant display font. Conveys a formal Gothic or supernatural theme.
- **Example Usage**: `style := Style{FontID: 10, EffectID: 3, Colors: []int{16777215}}` (Original flagship Hedwig lookup style).
- **Compatibility Status**: Extremely High.

### 11. GG Sans (Default)
- **ID**: `11`
- **Key Name**: `default`
- **Description**: Discord's standard native interface font. Acts as a reset / clean fallback profile display.
- **Example Usage**: `style := Style{FontID: 11, EffectID: 1, Colors: []int{16777215}}`
- **Compatibility Status**: Standard global default.

### 12. Zilla Slab (Tempo)
- **ID**: `12`
- **Key Name**: `zilla_slab`
- **Description**: Heavy, modern slab-serif typography optimized for reading at distance. Bold slab terminals with standard tracking.
- **Example Usage**: `style := Style{FontID: 12, EffectID: 1, Colors: []int{5865}}`
- **Compatibility Status**: High.

---

## ⏳ Future Support Notes
As Discord continues rolling out display metadata upgrades, some font names or IDs may receive localized adjustments. The style engine handles this by:
1. Conducting a startup **discovery test**.
2. Checking if selected font configurations throw validation `400` errors.
3. Automatically logging any mismatched font parameter names as `unsupportedFields` to protect the bot lifecycle.
4. 
