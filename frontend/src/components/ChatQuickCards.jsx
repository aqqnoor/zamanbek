// src/components/ChatQuickCards.jsx
export default function ChatQuickCards({ onQuickSelect }) {
  const options = [
    {
      title: "Ниет қою",
      description: "Craft compelling text for ads and emails."
    },
    {
      title: "Садақа жоспары",
      description: "Write articles on any topic instantly."
    },
    {
      title: "Кеңес алу",
      description: "Design custom visuals with AI."
    }
  ];

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 max-w-5xl mx-auto mt-6">
      {options.map((item, i) => (
        <button
          key={i}
          onClick={() => onQuickSelect(item.title)}
          className="p-5 bg-white rounded-3xl shadow-md hover:shadow-xl hover:scale-[1.03] transition text-left text-sm border border-gray-100"
        >
          <div className="font-semibold text-gray-900 mb-1 text-base">{item.title}</div>
          <div className="text-gray-500 text-sm leading-snug">
            {item.description}
          </div>
        </button>
      ))}
    </div>
  );
}
