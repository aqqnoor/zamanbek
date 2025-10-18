// src/components/ChatMessageInput.jsx
import { useState } from "react";

export default function ChatMessageInput({ onSend, loading, recording, startRecording, stopRecording }) {
  const [input, setInput] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!input.trim()) return;
    onSend(input);
    setInput("");
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="mt-4 bg-white px-4 py-2 rounded-full shadow-md flex items-center gap-3 max-w-3xl w-full mx-auto"
    >
      {/* Icons left side */}
      <div className="flex gap-2 text-gray-400 text-lg">
        <button type="button" disabled className="hover:text-gray-600">📎</button>
        <button
          type="button"
          onClick={recording ? stopRecording : startRecording}
          className="hover:text-gray-600"
        >
          {recording ? "⏹️" : "🎙️"}
        </button>
      </div>

      {/* Input */}
      <input
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Write your message ..."
        className="flex-1 bg-transparent outline-none text-sm"
      />

      {/* Send Button */}
      <button
        type="submit"
        disabled={loading}
        className="text-[#00A884] hover:text-green-600 text-xl transition"
      >
        ➤
      </button>
    </form>
  );
}