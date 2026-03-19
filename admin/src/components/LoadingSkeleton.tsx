export function CardSkeleton({ count = 3 }: { count?: number }) {
  return (
    <div className="space-y-3 animate-pulse">
      {Array.from({ length: count }).map((_, i) => (
        <div key={i} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <div className="h-4 bg-gray-800 rounded w-1/3 mb-3" />
          <div className="h-3 bg-gray-800 rounded w-2/3" />
        </div>
      ))}
    </div>
  )
}

export function TableSkeleton({ rows = 5, cols = 4 }: { rows?: number; cols?: number }) {
  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden animate-pulse">
      <div className="border-b border-gray-800 px-5 py-3 flex gap-8">
        {Array.from({ length: cols }).map((_, i) => (
          <div key={i} className="h-3 bg-gray-800 rounded w-20" />
        ))}
      </div>
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="px-5 py-4 flex gap-8 border-b border-gray-800/50">
          {Array.from({ length: cols }).map((_, j) => (
            <div key={j} className="h-3 bg-gray-800 rounded w-16" />
          ))}
        </div>
      ))}
    </div>
  )
}

export function StatSkeleton() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 animate-pulse">
      {[1, 2, 3].map((i) => (
        <div key={i} className="bg-gray-900 border border-gray-800 rounded-xl p-5">
          <div className="h-3 bg-gray-800 rounded w-20 mb-3" />
          <div className="h-6 bg-gray-800 rounded w-12" />
        </div>
      ))}
    </div>
  )
}
