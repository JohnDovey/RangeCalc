# RangeCalc GPS + Map

**A powerful offline Progressive Web App for calculating range to a target using GPS coordinates and compass bearings.**

![RangeCalc GPS](https://via.placeholder.com/800x400/0d6efd/ffffff?text=RangeCalc+GPS+%2B+Map)

---

## ✨ Features

- **GPS Integration** — Automatically capture location for Ref 1, Ref 2, and the Target
- **Dual Calculation Modes**:
  - Simple (fast)
  - Full Triangulation (Law of Sines)
- **Interactive Leaflet Map** — Visual verification with markers and range circles
- **Target GPS Verification** — Compare calculated range vs real measured distance
- **Baseline Auto-calculation** — Exact distance and bearing from GPS coordinates
- **Dark / Light Mode** — Toggleable theme
- **Fully Offline** — Works as an installable PWA (no internet required after first load)
- **Single File** — Easy to distribute and host

---

## How to Use

1. **Open the app** in Chrome (or any modern browser)
2. Tap **📍 Get GPS** for **Ref 1** and **Ref 2**
3. (Optional but recommended) Move to the actual target and tap **📍 Get Target GPS** for verification
4. Enter the two compass **bearings** to the target
5. Choose calculation method and tap **Calculate & Show Map**
6. View numerical results + interactive map showing:
   - Ref 1 & Ref 2 markers
   - Baseline line
   - Range circles from each reference point
   - Actual Target marker (if provided)

The intersection of the two range circles shows the calculated target location.

---

## Installation (PWA)

1. Open the HTML file in **Chrome** on your Android phone
2. Tap the **three-dot menu** → **"Add to Home screen"** (or **Install**)
3. The app will appear as a native-like icon and works **completely offline**

---

## Technical Details

- Pure HTML + JavaScript (single file)
- Uses **Leaflet.js** for the interactive map
- Haversine formula for GPS distance
- Law of Sines for triangulation
- No backend or data collection

---

## Author

**John Dovey**  
- Email: dovey.john@gmail.com  
- GitHub: [JohnDovey](https://github.com/JohnDovey)  
- Project: [RangeCalc](https://github.com/JohnDovey/RangeCalc)

---

## License

Open source and free to use / modify.  
If this tool helps you in the field, feel free to [donate via PayPal](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=C9QFERPZUK7D2).

---

**Made for surveyors, hunters, search & rescue, and outdoor professionals.**

---

*Last updated: May 2026*
