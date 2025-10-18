// src/hooks/useRecorder.js
import { useState, useRef } from "react";

export default function useRecorder(onTranscribe) {
  const [recording, setRecording] = useState(false);
  const mediaRecorderRef = useRef(null);
  const chunksRef = useRef([]);

  const startRecording = async () => {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    const mediaRecorder = new MediaRecorder(stream);
    mediaRecorderRef.current = mediaRecorder;
    chunksRef.current = [];

    mediaRecorder.ondataavailable = (e) => chunksRef.current.push(e.data);

    mediaRecorder.onstop = async () => {
      const blob = new Blob(chunksRef.current, { type: "audio/webm" });
      const file = new File([blob], "recording.webm", { type: "audio/webm" });
      const form = new FormData();
      form.append("file", file);
      form.append("model", "whisper-1");

      try {
        const res = await fetch("https://openai-hub.neuraldeep.tech/v1/audio/transcriptions", {
          method: "POST",
          headers: {
            Authorization: "Bearer sk-roG3OusRr0TLCHAADks6lw",
          },
          body: form,
        });
        const data = await res.json();
        const text = data.text;
        onTranscribe(text);
      } catch (err) {
        console.error("Transcription error", err);
      }
    };

    mediaRecorder.start();
    setRecording(true);
  };

  const stopRecording = () => {
    mediaRecorderRef.current?.stop();
    setRecording(false);
  };

  return { recording, startRecording, stopRecording };
}
