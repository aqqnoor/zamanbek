import { useNavigate } from "react-router-dom";
import FloatingChatButton from "../components/FloatingChatButton";
import heroDesktop from "../assets/hero-desktop.png";
import heroTablet from "../assets/hero-tablet.png";
import heroMobile from "../assets/hero-mobile.png";

export default function MainPage() {
  const navigate = useNavigate();

  return (
    <div className="relative w-full">
      <picture>
        <source srcSet={heroDesktop} media="(min-width:1024px)" />
        <source srcSet={heroTablet} media="(min-width:640px)" />

        {/* ВАЖНО: картинка в потоке, без absolute.
            Если файл длинный — появится естественный скролл. */}
        <img src={heroMobile} alt="background" className="w-full h-auto block" />
      </picture>

      <FloatingChatButton onClick={() => navigate("/chat")} />
    </div>
  );
}
