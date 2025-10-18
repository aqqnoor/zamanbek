// адрес твоего Go-бэкенда; если будешь делать прокси в Vite — поставь BASE = "/api"
const BASE = import.meta.env.VITE_API_BASE || "http://localhost:8080";

export async function sendMessage(message) {
  try {
    const res = await fetch(`${BASE}/chat`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ message }),
    });
    const data = await res.json();
    return data.reply || "Жауап жоқ.";
  } catch (e) {
    console.error(e);
    return "Қате пайда болды.";
  }
}
