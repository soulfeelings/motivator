import { useState, useEffect } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import { Book, ChevronRight } from 'lucide-react'

const docs = [
  { id: 'getting-started', title: 'Getting Started', file: '/docs/getting-started.md' },
  { id: 'admin-guide', title: 'Admin Guide', file: '/docs/admin-guide.md' },
  { id: 'api-reference', title: 'API Reference', file: '/docs/api-reference.md' },
  { id: 'integrations-guide', title: 'Integrations', file: '/docs/integrations-guide.md' },
  { id: 'game-server', title: 'Command Center', file: '/docs/game-server.md' },
]

export default function Docs() {
  const [activeDoc, setActiveDoc] = useState(docs[0].id)
  const [content, setContent] = useState('')
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadDoc(activeDoc)
  }, [activeDoc])

  async function loadDoc(id: string) {
    setLoading(true)
    const doc = docs.find(d => d.id === id)
    if (!doc) return
    try {
      const res = await fetch(doc.file)
      const text = await res.text()
      setContent(text)
    } catch {
      setContent('# Error\n\nFailed to load documentation.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex gap-6 -mx-8 -my-8 h-[calc(100vh-0px)]">
      {/* Sidebar */}
      <div className="w-64 border-r border-gray-800 bg-gray-950/50 p-5 shrink-0">
        <div className="flex items-center gap-2 mb-6">
          <Book size={20} className="text-violet-400" />
          <h2 className="text-lg font-bold text-white">Docs</h2>
        </div>
        <nav className="space-y-1">
          {docs.map(doc => (
            <button
              key={doc.id}
              onClick={() => setActiveDoc(doc.id)}
              className={`w-full flex items-center gap-2 px-3 py-2.5 rounded-lg text-sm text-left transition-colors ${
                activeDoc === doc.id
                  ? 'bg-violet-600/20 text-violet-400 font-medium'
                  : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/50'
              }`}
            >
              <ChevronRight size={14} className={activeDoc === doc.id ? 'text-violet-400' : 'text-gray-600'} />
              {doc.title}
            </button>
          ))}
        </nav>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-8">
        {loading ? (
          <div className="animate-pulse space-y-4">
            <div className="h-8 bg-gray-800 rounded w-1/3" />
            <div className="h-4 bg-gray-800 rounded w-2/3" />
            <div className="h-4 bg-gray-800 rounded w-1/2" />
          </div>
        ) : (
          <div className="prose prose-invert prose-violet max-w-3xl
            prose-headings:text-white prose-headings:font-bold
            prose-h1:text-3xl prose-h1:border-b prose-h1:border-gray-800 prose-h1:pb-4 prose-h1:mb-6
            prose-h2:text-xl prose-h2:mt-8 prose-h2:mb-4
            prose-h3:text-lg prose-h3:mt-6 prose-h3:mb-3
            prose-p:text-gray-400 prose-p:leading-relaxed
            prose-a:text-violet-400 prose-a:no-underline hover:prose-a:underline
            prose-strong:text-white
            prose-code:text-violet-400 prose-code:bg-gray-800/50 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:text-sm prose-code:before:content-none prose-code:after:content-none
            prose-pre:bg-gray-900 prose-pre:border prose-pre:border-gray-800 prose-pre:rounded-xl
            prose-table:border-collapse
            prose-th:bg-gray-900 prose-th:px-4 prose-th:py-2 prose-th:text-left prose-th:text-xs prose-th:uppercase prose-th:text-gray-500 prose-th:border-b prose-th:border-gray-800
            prose-td:px-4 prose-td:py-2.5 prose-td:text-sm prose-td:text-gray-400 prose-td:border-b prose-td:border-gray-800/50
            prose-li:text-gray-400
            prose-hr:border-gray-800
          ">
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{content}</ReactMarkdown>
          </div>
        )}
      </div>
    </div>
  )
}
