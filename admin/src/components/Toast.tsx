import { useEffect, useState } from 'react'
import { X, CheckCircle, AlertCircle } from 'lucide-react'

interface ToastProps {
  message: string
  type?: 'success' | 'error'
  onClose: () => void
}

export default function Toast({ message, type = 'success', onClose }: ToastProps) {
  const [visible, setVisible] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setVisible(false)
      setTimeout(onClose, 300)
    }, 4000)
    return () => clearTimeout(timer)
  }, [onClose])

  const colors = type === 'success'
    ? 'bg-emerald-900/90 border-emerald-700 text-emerald-200'
    : 'bg-red-900/90 border-red-700 text-red-200'

  return (
    <div className={`fixed bottom-6 right-6 z-50 flex items-center gap-3 px-4 py-3 rounded-xl border shadow-xl transition-all duration-300 ${colors} ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
      {type === 'success' ? <CheckCircle size={18} /> : <AlertCircle size={18} />}
      <span className="text-sm font-medium">{message}</span>
      <button onClick={() => { setVisible(false); setTimeout(onClose, 300) }} className="ml-2 opacity-60 hover:opacity-100">
        <X size={14} />
      </button>
    </div>
  )
}
