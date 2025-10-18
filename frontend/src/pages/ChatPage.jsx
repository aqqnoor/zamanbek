// src/pages/ChatPage.jsx
import { useState } from "react";
import ChatSidebar from "../components/ChatSidebar";
import ChatStatsSidebar from "../components/ChatStatsSidebar";
import ChatQuickCards from "../components/ChatQuickCards";
import ChatMessageInput from "../components/ChatMessageInput";
import useRecorder from "../hooks/useRecorder";

export default function ChatPage() {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(false);

  const sendMessage = async (text) => {
    if (!text.trim()) return;
    const newUserMsg = { role: "user", content: text };
    setMessages((prev) => [...prev, newUserMsg]);
    setLoading(true);

    try {
      const res = await fetch("https://openai-hub.neuraldeep.tech/v1/chat/completions", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer sk-roG3OusRr0TLCHAADks6lw",
        },
        body: JSON.stringify({
          model: "gpt-4o-mini",
          messages: [...messages, newUserMsg],
        }),
      });
      const data = await res.json();
      const reply = data.choices?.[0]?.message?.content || "...";
      setMessages((prev) => [...prev, { role: "assistant", content: reply }]);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const { recording, startRecording, stopRecording } = useRecorder((text) => {
    sendMessage(text);
  });

  return (
    <div className="flex h-svh overflow-hidden font-sans">
      <ChatSidebar />

      <main className="flex-1 flex flex-col px-6 py-6 overflow-y-auto bg-gradient-to-br from-[#E6F1EF] to-[#C8E7D4]">
        <h1 className="text-3xl sm:text-4xl font-bold text-center text-[#00A884]">
          Заманбекке қош келдіңіз!
        </h1>
        <p className="text-center text-gray-600 text-sm mb-6">
          Start by scripting a task, and let the chat take over.
        </p>

        <div className="mx-auto w-20 h-20 rounded-full bg-white shadow-md flex items-center justify-center mb-6">
          🤖
        </div>

        <ChatQuickCards onQuickSelect={sendMessage} />

        <div className="flex flex-col gap-2 mt-6 mb-4 overflow-y-auto flex-1">
          {messages.map((msg, i) => (
            <div
              key={i}
              className={`max-w-[80%] px-4 py-2 rounded-xl text-sm ${
                msg.role === "user"
                  ? "self-end bg-green-100"
                  : "self-start bg-white"
              }`}
            >
              {msg.content}
            </div>
          ))}
          {loading && (
            <div className="self-start bg-white text-gray-500 px-4 py-2 rounded-xl text-sm">
              ...
            </div>
          )}
        </div>

        <ChatMessageInput
          onSend={sendMessage}
          loading={loading}
          recording={recording}
          startRecording={startRecording}
          stopRecording={stopRecording}
        />
      </main>

      <ChatStatsSidebar />
    </div>
  );
}
