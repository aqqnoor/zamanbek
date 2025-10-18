import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { sendMessage } from "../services/api";

export default function ChatPage() {
  const navigate = useNavigate();
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState("");

  const handleSend = async () => {
    if (!input.trim()) return;
    const userMsg = { role: "user", text: input };
    setMessages(prev => [...prev, userMsg]);
    setInput("");
    const reply = await sendMessage(input);
    setMessages(prev => [...prev, { role: "assistant", text: reply }]);
  };

  return (
    <div className="h-screen flex flex-col bg-gray-50">
      <header className="flex items-center gap-3 p-4 bg-white shadow">
        <button onClick={() => navigate("/")} className="text-[#00A884]">← Назад</button>
        <h2 className="text-lg font-semibold">Ассистент Zaman</h2>
      </header>

      <main className="flex-1 overflow-auto p-4 space-y-2">
        {messages.map((m, i) => (
          <div key={i}
            className={`max-w-[80%] p-3 rounded-xl ${m.role === "user" ? "ml-auto bg-green-100" : "mr-auto bg-white border"}`}>
            {m.text}
          </div>
        ))}
      </main>

      <footer className="p-3 bg-white border-t flex gap-2">
        <input
          className="flex-1 border rounded-lg px-3 py-2"
          placeholder="Жазыңыз..."
          value={input}
          onChange={(e)=>setInput(e.target.value)}
          onKeyDown={(e)=>e.key==='Enter' && handleSend()}
        />
        <button onClick={handleSend} className="bg-[#00A884] text-white px-4 py-2 rounded-lg">→</button>
      </footer>
    </div>
  );
}
