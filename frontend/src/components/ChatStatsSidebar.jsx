// ChatStatsSidebar.jsx
export default function ChatStatsSidebar() {
  return (
    <aside className="w-64 bg-gray-50 p-4 flex flex-col space-y-4">
      {/* Карточка Уровень Береке */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col items-center">
        <h3 className="text-lg font-semibold mb-2">Береке</h3>
        {/* Пример кругового прогресса (можно заменить на реальный компонент) */}
        <div className="relative">
          <svg className="w-20 h-20" viewBox="0 0 100 100">
            <circle cx="50" cy="50" r="45" strokeWidth="10" fill="transparent" 
                    className="text-gray-200 stroke-current"/>
            <circle cx="50" cy="50" r="45" strokeWidth="10" fill="transparent" 
                    className="text-indigo-500 stroke-current 
                               transform -rotate-90 origin-center"
                    strokeDasharray="282.6" strokeDashoffset="84.78" 
                    strokeLinecap="round" />
            <text x="50" y="50" textAnchor="middle" dy="0.4em" 
                  className="text-sm fill-current text-gray-800">70%</text>
          </svg>
        </div>
      </div>

      {/* Карточка Күндік лимит */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col items-start">
        <h3 className="text-lg font-semibold mb-2">Күндік лимит</h3>
        {/* Полоса прогресса дневного лимита */}
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div className="bg-green-500 h-2 rounded-full" style={{ width: '50%' }}></div>
        </div>
        <span className="text-sm text-gray-600 mt-1">50% пайдаланылды</span>
      </div>

      {/* Карточка Бүгінгі миссиясы */}
      <div className="bg-white rounded-xl shadow p-4 flex flex-col items-start">
        <h3 className="text-lg font-semibold mb-2">Бүгінгі миссиясы</h3>
        {/* Мини-график (пример столбиками) */}
        <div className="flex items-end space-x-1 h-12 w-full bg-gray-100 rounded px-2 py-1">
          <div className="bg-blue-500 w-2 h-8 rounded-t"></div>
          <div className="bg-blue-500 w-2 h-5 rounded-t"></div>
          <div className="bg-blue-500 w-2 h-9 rounded-t"></div>
          <div className="bg-blue-500 w-2 h-4 rounded-t"></div>
          <div className="bg-blue-500 w-2 h-7 rounded-t"></div>
        </div>
        <span className="text-sm text-gray-600 mt-1">8/10 аяқталды</span>
      </div>
    </aside>
  );
}
