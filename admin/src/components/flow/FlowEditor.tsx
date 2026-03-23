import { useCallback, useRef, useState, useEffect } from 'react'
import {
  ReactFlow,
  addEdge,
  useNodesState,
  useEdgesState,
  Controls,
  Background,
  BackgroundVariant,
} from '@xyflow/react'
import type { Connection, ReactFlowInstance, Node, Edge, NodeChange } from '@xyflow/react'
import '@xyflow/react/dist/style.css'

import TriggerNode from './TriggerNode'
import ConditionNode from './ConditionNode'
import ActionNode from './ActionNode'
import NodeSidebar from './NodeSidebar'
import NodeEditPanel from './NodeEditPanel'

const nodeTypes = {
  trigger: TriggerNode,
  condition: ConditionNode,
  action: ActionNode,
}

const defaultNodeData: Record<string, Record<string, unknown>> = {
  trigger: { metric: 'deals_closed', description: 'When metric is reported' },
  condition: { metric: 'deals_closed', operator: 'gte', threshold: 10 },
  action: { label: 'Award Reward', xp: 50, coins: 25, badge: false },
}

interface FlowEditorProps {
  initialNodes?: Node[]
  initialEdges?: Edge[]
  onSave: (nodes: Node[], edges: Edge[]) => void
  saving?: boolean
}

let nodeId = 0
function getNodeId() {
  return `node_${++nodeId}_${Date.now()}`
}

export default function FlowEditor({ initialNodes = [], initialEdges = [], onSave, saving }: FlowEditorProps) {
  const reactFlowWrapper = useRef<HTMLDivElement>(null)
  const [rfInstance, setRfInstance] = useState<ReactFlowInstance | null>(null)
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)

  const onNodeClick = useCallback((_event: React.MouseEvent, node: Node) => {
    setSelectedNode(node)
  }, [])

  const onPaneClick = useCallback(() => {
    setSelectedNode(null)
  }, [])

  const onUpdateNodeData = useCallback((nodeId: string, data: Record<string, unknown>) => {
    setNodes((nds) => nds.map((n) => n.id === nodeId ? { ...n, data } : n))
    setSelectedNode((prev) => prev && prev.id === nodeId ? { ...prev, data } : prev)
  }, [setNodes])

  const handleNodesChange = useCallback((changes: NodeChange[]) => {
    onNodesChange(changes)
    const removals = changes.filter((c): c is NodeChange & { type: 'remove'; id: string } => c.type === 'remove')
    if (selectedNode && removals.some((r) => r.id === selectedNode.id)) {
      setSelectedNode(null)
    }
  }, [onNodesChange, selectedNode])

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setSelectedNode(null)
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge({ ...params, animated: true, style: { stroke: '#6b7280' } }, eds)),
    [setEdges],
  )

  const onDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'move'
  }, [])

  const onDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault()
      const type = e.dataTransfer.getData('application/reactflow')
      if (!type || !rfInstance || !reactFlowWrapper.current) return

      const bounds = reactFlowWrapper.current.getBoundingClientRect()
      const position = rfInstance.screenToFlowPosition({
        x: e.clientX - bounds.left,
        y: e.clientY - bounds.top,
      })

      const newNode: Node = {
        id: getNodeId(),
        type,
        position,
        data: { ...defaultNodeData[type] },
      }

      setNodes((nds) => [...nds, newNode])
    },
    [rfInstance, setNodes],
  )

  return (
    <div className="flex h-[600px] border border-gray-800 rounded-xl overflow-hidden">
      <NodeSidebar />
      <div className="flex-1 relative" ref={reactFlowWrapper}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={handleNodesChange}
          onEdgesChange={onEdgesChange}
          onConnect={onConnect}
          onInit={setRfInstance}
          onDrop={onDrop}
          onDragOver={onDragOver}
          onNodeClick={onNodeClick}
          onPaneClick={onPaneClick}
          nodeTypes={nodeTypes}
          fitView
          colorMode="dark"
          defaultEdgeOptions={{ animated: true, style: { stroke: '#6b7280' } }}
        >
          <Controls className="!bg-gray-800 !border-gray-700 !rounded-lg [&>button]:!bg-gray-800 [&>button]:!border-gray-700 [&>button]:!text-gray-400 [&>button:hover]:!bg-gray-700" />
          <Background variant={BackgroundVariant.Dots} gap={20} size={1} color="#374151" />
        </ReactFlow>
        <div className="absolute top-4 right-4 z-10">
          <button
            onClick={() => onSave(nodes, edges)}
            disabled={saving}
            className="px-4 py-2 bg-violet-600 hover:bg-violet-500 disabled:opacity-50 text-white text-sm font-medium rounded-lg transition-colors shadow-lg"
          >
            {saving ? 'Saving...' : 'Save Flow'}
          </button>
        </div>
      </div>
      {selectedNode && (
        <NodeEditPanel
          node={selectedNode}
          onUpdate={onUpdateNodeData}
          onClose={() => setSelectedNode(null)}
        />
      )}
    </div>
  )
}
