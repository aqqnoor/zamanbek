export default function ChatSidebar() {
const items = ["Chat", "Habits", "Goals", "Rewards"];
return (
<aside className="hidden md:flex flex-col w-48 bg-white border-r p-4 text-sm text-gray-700">
{items.map((item, i) => (
<div
key={i}
className={`mb-2 px-3 py-2 rounded-lg transition-all cursor-pointer hover:bg-gray-100 ${
item === "Chat" ? "bg-[#00A884] text-white" : ""
}`}
>
{item}
</div>
))}
</aside>
);
}