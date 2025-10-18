import chatIcon from "../assets/chat-ball.png";

export default function FloatingChatButton({ onClick }) {
  return (
    <button
      onClick={onClick}
      aria-label="Open chat"
      className="fixed z-50 bottom-6 right-6 w-16 h-16 rounded-full flex items-center justify-center
                 hover:scale-110 transition-transform duration-300"
      style={{
        background: "transparent",
        border: "none",
        boxShadow: "0 8px 25px rgba(0,0,0,.15)",
      }}
    >
      {/* Кольцо звёзд поверх кнопки */}
      <div className="sparkle-ring">
        <div className="sparkle-orbit">
          {/* 8 звёзд по окружности. Можно увеличить/уменьшить количество, шаг меняется в CSS (--deg). */}
          <span className="sparkle" style={{'--i': 0}}></span>
          <span className="sparkle" style={{'--i': 1}}></span>
          <span className="sparkle" style={{'--i': 2}}></span>
          <span className="sparkle" style={{'--i': 3}}></span>
          <span className="sparkle" style={{'--i': 4}}></span>
          <span className="sparkle" style={{'--i': 5}}></span>
          <span className="sparkle" style={{'--i': 6}}></span>
          <span className="sparkle" style={{'--i': 7}}></span>
        </div>
      </div>

      {/* Сам шарик-иконка */}
      <img
        src={chatIcon}
        alt="Chat"
        className="w-16 h-16 object-contain"
        style={{ filter: "drop-shadow(0 4px 10px rgba(0,0,0,.2))" }}
      />
    </button>
  );
}
