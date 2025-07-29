export const THEMES = [
  "light",
  "dark",
  "cupcake",
  "bumblebee",
  "emerald",
  "corporate",
  "synthwave",
  "retro",
  "cyberpunk",
  "valentine",
  "halloween",
  "garden",
  "forest",
  "aqua",
  "lofi",
  "pastel",
  "fantasy",
  "wireframe",
  "black",
  "luxury",
  "dracula",
  "cmyk",
  "autumn",
  "business",
  "acid",
  "lemonade",
  "night",
  "coffee",
  "winter",
  "dim",
  "nord",
  "sunset",
] as const; // List of available themes


export const API_BASE_URL = "http://localhost:3000/api";
export const PROD_API_BASE_URL = "https://your-production-api.com/api"; // Replace with your production API URL
export const BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:5001" : "/";