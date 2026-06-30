import React, { useEffect, useState, useCallback } from 'react'
import {
  ChevronRight, Plus, Trash2, Pencil, Zap, Wrench,
  CheckCircle, XCircle, Clock, ArrowLeft, X,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import {
  ListWorkflows, ListSources, GetWorkflow,
  SaveWorkflow, DeleteWorkflow, GetWorkflowDiff, RunWorkflow,
} from '../../wailsjs/go/gui/App'
import type { Workflow, Step, Endpoint, Source, StepDiff, StepResult, WorkflowResult } from '@/types'

// ── helpers ────────────────────────────────────────────────────────────────

function kindIcon(kind: string) {
  switch (kind.toLowerCase()) {
    case 'sync': return <Zap    className="h-3.5 w-3.5" />
    case 'fix':  return <Wrench className="h-3.5 w-3.5" />
    default:     return <Pencil className="h-3.5 w-3.5" />
  }
}
function kindVariant(kind: string): 'sync' | 'fix' | 'edit' {
  if (kind === 'sync') return 'sync'
  if (kind === 'fix')  return 'fix'
  return 'edit'
}
function statusIcon(status: string) {
  if (status === 'success') return <CheckCircle className="h-4 w-4 text-emerald-400" />
  if (status === 'failed')  return <XCircle     className="h-4 w-4 text-red-400"     />
  return <Clock className="h-4 w-4 text-purple-400" />
}

// cast helpers — Wails returns class instances; we use plain interfaces in state
const asWorkflows = (x: unknown) => (x ?? []) as Workflow[]
const asSources   = (x: unknown) => (x ?? []) as Source[]
const asWorkflow  = (x: unknown) => x as Workflow
const asDiffs     = (x: unknown) => (x ?? []) as StepDiff[]
const asResult    = (x: unknown) => x as WorkflowResult

function blankStep(srcId: string): Step {
  return {
    id: '', kind: 'sync',
    source:  { source_id: srcId, resource: 'tracks',    query: '' },
    targets: [{ source_id: srcId, resource: 'playlists', query: '' }],
    after: [], options: {},
  }
}

// ── WorkflowsView ──────────────────────────────────────────────────────────

type Mode = 'list' | 'view' | 'edit' | 'applying'

export default function WorkflowsView() {
  const [mode,     setMode]     = useState<Mode>('list')
  const [wfList,   setWfList]   = useState<Workflow[]>([])
  const [selected, setSelected] = useState<Workflow | null>(null)
  const [sources,  setSources]  = useState<Source[]>([])
  const [diffs,    setDiffs]    = useState<StepDiff[]>([])
  const [result,   setResult]   = useState<WorkflowResult | null>(null)
  const [error,    setError]    = useState('')
  const [busy,     setBusy]     = useState(false)

  const load = useCallback(async () => {
    try {
      const [wfs, srcs] = await Promise.all([ListWorkflows(), ListSources()])
      setWfList(asWorkflows(wfs)); setSources(asSources(srcs))
    } catch (e) { setError(String(e)) }
  }, [])

  useEffect(() => { load() }, [load])

  async function fetchDiff(id: string) {
    setBusy(true); setError('')
    try { setDiffs(asDiffs(await GetWorkflowDiff(id))) }
    catch (e) { setError(String(e)) }
    setBusy(false)
  }

  async function openWorkflow(w: Workflow) {
    setError(''); setDiffs([]); setResult(null)
    try {
      const full = asWorkflow(await GetWorkflow(w.id))
      setSelected(JSON.parse(JSON.stringify(full)))
      setMode('view'); fetchDiff(full.id)
    } catch (e) { setError(String(e)) }
  }

  async function handleNew() {
    setBusy(true); setError('')
    try {
      const wf = asWorkflow(await SaveWorkflow({ id: '', name: 'New Workflow', steps: [] } as never))
      await load(); setSelected(wf); setDiffs([]); setResult(null); setMode('edit')
    } catch (e) { setError(String(e)) }
    setBusy(false)
  }

  async function handleDelete(id: string, name: string) {
    if (!confirm(`Delete "${name}"?`)) return
    try { await DeleteWorkflow(id); await load() }
    catch (e) { setError(String(e)) }
  }

  async function handleSave(wf: Workflow) {
    setBusy(true); setError('')
    try {
      const saved = asWorkflow(await SaveWorkflow(wf as never))
      setSelected(saved); await load(); setMode('view'); fetchDiff(saved.id)
    } catch (e) { setError(String(e)) }
    setBusy(false)
  }

  async function handleApply() {
    if (!selected || !confirm('Apply all changes? This cannot be undone.')) return
    setBusy(true); setError(''); setMode('applying')
    try { setResult(asResult(await RunWorkflow(selected.id))) }
    catch (e) { setError(String(e)) }
    setBusy(false)
  }

  function backToList() { setMode('list'); setSelected(null); setDiffs([]); setResult(null); setError('') }

  // ── LIST ──────────────────────────────────────────────────────────────────
  if (mode === 'list') return (
    <div className="p-6 max-w-2xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-lg font-semibold">Workflows</h1>
          <p className="text-sm text-muted-foreground mt-0.5">Automate your library maintenance</p>
        </div>
        <Button size="sm" onClick={handleNew} disabled={busy}>
          <Plus className="h-4 w-4 mr-1.5" /> New Workflow
        </Button>
      </div>
      {error && <p className="text-sm text-destructive mb-4">{error}</p>}
      {wfList.length === 0
        ? <p className="text-sm text-muted-foreground italic">No workflows yet.</p>
        : <div className="flex flex-col gap-2">
            {wfList.map(w => (
              <Card key={w.id} className="cursor-pointer hover:border-border/80 transition-colors"
                    onClick={() => openWorkflow(w)}>
                <CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
                  <div className="flex items-center gap-3">
                    <span className="text-sm font-medium">{w.name}</span>
                    <span className="text-xs text-muted-foreground">
                      {w.steps?.length ?? 0} step{w.steps?.length !== 1 ? 's' : ''}
                    </span>
                  </div>
                  <div className="flex items-center gap-1">
                    <Button variant="ghost" size="icon"
                      onClick={e => { e.stopPropagation(); handleDelete(w.id, w.name) }}>
                      <Trash2 className="h-3.5 w-3.5 text-muted-foreground" />
                    </Button>
                    <ChevronRight className="h-4 w-4 text-muted-foreground" />
                  </div>
                </CardHeader>
              </Card>
            ))}
          </div>
      }
    </div>
  )

  // ── EDIT ──────────────────────────────────────────────────────────────────
  if (mode === 'edit' && selected) return (
    <WorkflowEditor
      workflow={selected}
      sources={sources}
      busy={busy}
      error={error}
      onSave={handleSave}
      onCancel={() => selected.id ? (setMode('view'), fetchDiff(selected.id)) : backToList()}
    />
  )

  // ── VIEW / APPLYING ────────────────────────────────────────────────────────
  if ((mode === 'view' || mode === 'applying') && selected) return (
    <WorkflowDetail
      workflow={selected}
      diffs={diffs}
      result={result}
      mode={mode}
      busy={busy}
      error={error}
      onEdit={() => setMode('edit')}
      onApply={handleApply}
      onPreviewAgain={() => fetchDiff(selected.id)}
      onBack={backToList}
    />
  )

  return null
}

// ── WorkflowDetail ─────────────────────────────────────────────────────────

interface DetailProps {
  workflow:       Workflow
  diffs:          StepDiff[]
  result:         WorkflowResult | null
  mode:           'view' | 'applying'
  busy:           boolean
  error:          string
  onEdit:         () => void
  onApply:        () => void
  onPreviewAgain: () => void
  onBack:         () => void
}

function WorkflowDetail({ workflow, diffs, result, mode, busy, error, onEdit, onApply, onPreviewAgain, onBack }: DetailProps) {
  const diffById: Record<string, StepDiff>   = Object.fromEntries(diffs.map(d => [d.step_id, d]))
  const resultById: Record<string, StepResult> = Object.fromEntries((result?.steps ?? []).map(r => [r.step_id, r]))
  const syncSteps = workflow.steps.filter(s => s.kind === 'sync').length
  const diffLoaded = diffs.length > 0 || syncSteps === 0

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
        <Button variant="ghost" size="sm" onClick={onBack}>
          <ArrowLeft className="h-4 w-4 mr-1.5" /> Workflows
        </Button>
        <Separator orientation="vertical" className="h-5 mx-1" />
        <span className="text-sm font-semibold">{workflow.name}</span>
        <div className="flex-1" />
        {error && <span className="text-xs text-destructive mr-2 max-w-xs truncate">{error}</span>}
        {busy  && <span className="text-xs text-muted-foreground mr-2">Loading…</span>}
        <Button variant="outline" size="sm" onClick={onEdit} disabled={busy}>Edit</Button>
        {mode === 'applying' && result
          ? <Button variant="outline" size="sm" onClick={onPreviewAgain}>Preview Again</Button>
          : <Button size="sm" onClick={onApply} disabled={!diffLoaded || busy}>▶ Apply</Button>
        }
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="flex flex-col gap-4 max-w-3xl">
          {workflow.steps.length === 0 && (
            <p className="text-sm text-muted-foreground italic">No steps. Press Edit to add some.</p>
          )}
          {workflow.steps.map((step, i) => (
            <StepViewCard
              key={step.id || i}
              step={step}
              index={i}
              diff={diffById[step.id]}
              result={resultById[step.id]}
              showResult={mode === 'applying'}
            />
          ))}
        </div>
      </div>
    </div>
  )
}

// ── StepViewCard ───────────────────────────────────────────────────────────

function StepViewCard({ step, index, diff, result, showResult }:
  { step: Step; index: number; diff?: StepDiff; result?: StepResult; showResult: boolean }) {

  const [showUnchanged, setShowUnchanged] = useState(false)
  const removedSet  = new Set(diff?.removed.map(t => t.id) ?? [])
  const unchanged   = (diff?.current ?? []).filter(t => !removedSet.has(t.id))

  return (
    <Card className={
      result?.status === 'success' ? 'border-emerald-900' :
      result?.status === 'failed'  ? 'border-red-900'     :
      result?.status === 'blocked' ? 'border-purple-900'  : ''
    }>
      <CardHeader className="bg-[hsl(240_10%_6%)] rounded-t-xl border-b border-border py-3 px-4">
        <div className="flex flex-wrap items-center gap-3">
          <span className="flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground shrink-0">
            {index + 1}
          </span>
          <Badge variant={kindVariant(step.kind)} className="flex items-center gap-1">
            {kindIcon(step.kind)} {step.kind.toUpperCase()}
          </Badge>
          <div className="flex items-center gap-1.5 min-w-0 flex-1 text-sm flex-wrap">
            <EndpointChip ep={step.source} />
            {step.targets.map((tgt, ti) => (
              <React.Fragment key={ti}>
                <ChevronRight className="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                <EndpointChip ep={tgt} />
              </React.Fragment>
            ))}
          </div>
          {showResult && result && (
            <div className="flex items-center gap-1.5 shrink-0">
              {statusIcon(result.status)}
              <span className="text-xs text-muted-foreground capitalize">{result.status}</span>
            </div>
          )}
        </div>
      </CardHeader>

      {diff && step.kind === 'sync' && (
        <CardContent className="pt-3 pb-3">
          {diff.added.length === 0 && diff.removed.length === 0 ? (
            <div className="flex items-center gap-2 text-xs text-emerald-400 bg-emerald-950/40 rounded-md px-3 py-2">
              <CheckCircle className="h-3.5 w-3.5" /> Already up to date
            </div>
          ) : (
            <TrackDiffTable
              target={diff.target_name}
              added={diff.added}
              removed={diff.removed}
              unchanged={unchanged}
              showUnchanged={showUnchanged}
              onToggleUnchanged={() => setShowUnchanged(v => !v)}
            />
          )}
          {result?.error && <p className="text-xs text-destructive mt-2">✗ {result.error}</p>}
        </CardContent>
      )}
    </Card>
  )
}

// ── EndpointChip ───────────────────────────────────────────────────────────

function EndpointChip({ ep }: { ep: Endpoint }) {
  return (
    <span className="truncate text-xs font-mono text-muted-foreground bg-muted/40 px-1.5 py-0.5 rounded">
      {ep.resource}{ep.query ? ' · ' + ep.query : ''}
    </span>
  )
}

// ── TrackDiffTable ─────────────────────────────────────────────────────────

import type { TrackRow } from '@/types'

function TrackDiffTable({ target, added, removed, unchanged, showUnchanged, onToggleUnchanged }:
  { target: string; added: TrackRow[]; removed: TrackRow[]; unchanged: TrackRow[];
    showUnchanged: boolean; onToggleUnchanged: () => void }) {

  type RowKind = 'add' | 'remove' | 'unchanged'
  type DisplayRow = { kind: RowKind; track: TrackRow }

  const rows: DisplayRow[] = [
    ...added.map(t   => ({ kind: 'add'       as RowKind, track: t })),
    ...removed.map(t => ({ kind: 'remove'    as RowKind, track: t })),
    ...(showUnchanged ? unchanged.map(t => ({ kind: 'unchanged' as RowKind, track: t })) : []),
  ]

  return (
    <div>
      {target && <p className="text-xs text-muted-foreground mb-2 font-medium">{target}</p>}
      <div className="rounded-md border border-border overflow-hidden">
        <table className="w-full text-xs">
          <thead>
            <tr className="border-b border-border bg-muted/20">
              <th className="w-5 py-1.5 pl-2" />
              <th className="py-1.5 px-2 text-left font-medium text-muted-foreground">Title</th>
              <th className="py-1.5 px-2 text-left font-medium text-muted-foreground hidden sm:table-cell">Artist</th>
              <th className="py-1.5 px-2 text-right font-medium text-muted-foreground w-12">BPM</th>
            </tr>
          </thead>
          <tbody>
            {rows.map(({ kind, track }) => (
              <tr key={track.id + kind}
                className={
                  kind === 'add'    ? 'border-l-2 border-l-emerald-500 bg-emerald-950/30' :
                  kind === 'remove' ? 'border-l-2 border-l-red-500     bg-red-950/30'     :
                  'border-l-2 border-l-transparent opacity-50'
                }>
                <td className="pl-2 text-center font-mono font-bold">
                  {kind === 'add'    ? <span className="text-emerald-400">+</span>
                 : kind === 'remove' ? <span className="text-red-400">−</span>
                 : <span className="text-muted-foreground">·</span>}
                </td>
                <td className="py-1.5 px-2 max-w-0"><span className="truncate block">{track.title || track.id}</span></td>
                <td className="py-1.5 px-2 text-muted-foreground hidden sm:table-cell max-w-0">
                  <span className="truncate block">{track.artist}</span>
                </td>
                <td className="py-1.5 px-2 text-right text-muted-foreground font-mono">{track.bpm}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      {unchanged.length > 0 && (
        <button onClick={onToggleUnchanged}
          className="mt-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors">
          {showUnchanged ? '↑ Hide unchanged' : `↓ Show ${unchanged.length} unchanged`}
        </button>
      )}
    </div>
  )
}

// ── WorkflowEditor ─────────────────────────────────────────────────────────

function WorkflowEditor({ workflow, sources, busy, error, onSave, onCancel }:
  { workflow: Workflow; sources: Source[]; busy: boolean; error: string;
    onSave: (wf: Workflow) => void; onCancel: () => void }) {

  const [wf, setWf] = useState<Workflow>(() => JSON.parse(JSON.stringify(workflow)))
  const firstSrcId = sources[0]?.id ?? ''

  const mutSteps = (fn: (steps: Step[]) => Step[]) =>
    setWf(w => ({ ...w, steps: fn([...w.steps]) }))

  const updStep = (i: number, patch: Partial<Step>) =>
    mutSteps(ss => { ss[i] = { ...ss[i], ...patch }; return ss })

  const updSource = (si: number, patch: Partial<Endpoint>) =>
    mutSteps(ss => { ss[si] = { ...ss[si], source: { ...ss[si].source, ...patch } }; return ss })

  const updTarget = (si: number, ti: number, patch: Partial<Endpoint>) =>
    mutSteps(ss => {
      const tgts = [...ss[si].targets]; tgts[ti] = { ...tgts[ti], ...patch }
      ss[si] = { ...ss[si], targets: tgts }; return ss
    })

  const addTarget = (si: number) =>
    mutSteps(ss => { ss[si] = { ...ss[si], targets: [...ss[si].targets, { source_id: firstSrcId, resource: 'playlists', query: '' }] }; return ss })

  const removeTarget = (si: number, ti: number) =>
    mutSteps(ss => { ss[si] = { ...ss[si], targets: ss[si].targets.filter((_, j) => j !== ti) }; return ss })

  return (
    <div className="flex flex-col h-full">
      <div className="flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
        <Button variant="ghost" size="sm" onClick={onCancel} disabled={busy}>
          <X className="h-4 w-4 mr-1" /> Cancel
        </Button>
        <Separator orientation="vertical" className="h-5 mx-1" />
        <input
          className="bg-transparent border-none text-sm font-semibold focus:outline-none w-64"
          value={wf.name}
          onChange={e => setWf(w => ({ ...w, name: e.target.value }))}
          placeholder="Workflow name"
        />
        <div className="flex-1" />
        {error && <span className="text-xs text-destructive mr-2 max-w-xs truncate">{error}</span>}
        <Button size="sm" onClick={() => onSave(wf)} disabled={busy}>
          {busy ? 'Saving…' : 'Save'}
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <div className="flex flex-col gap-3 max-w-3xl">
          {wf.steps.length === 0 && (
            <p className="text-sm text-muted-foreground italic py-2">No steps yet — add one below.</p>
          )}
          {wf.steps.map((step, si) => (
            <Card key={si} className="border-border/60">
              <CardHeader className="bg-[hsl(240_10%_6%)] rounded-t-xl border-b border-border py-2.5 px-4">
                <div className="flex items-center gap-3">
                  <span className="flex h-6 w-6 items-center justify-center rounded-full bg-muted text-xs font-bold text-muted-foreground shrink-0">
                    {si + 1}
                  </span>
                  <Select value={step.kind} onValueChange={k => updStep(si, { kind: k })}>
                    <SelectTrigger className="w-24 h-7"><SelectValue /></SelectTrigger>
                    <SelectContent>
                      <SelectItem value="sync">SYNC</SelectItem>
                      <SelectItem value="fix">FIX</SelectItem>
                      <SelectItem value="edit">EDIT</SelectItem>
                    </SelectContent>
                  </Select>
                  <div className="flex-1" />
                  <Button variant="ghost" size="icon" className="h-7 w-7"
                    onClick={() => mutSteps(ss => ss.filter((_, j) => j !== si))}>
                    <X className="h-3.5 w-3.5 text-muted-foreground" />
                  </Button>
                </div>
              </CardHeader>

              <CardContent className="pt-3 pb-4 flex flex-col gap-3">
                <div>
                  <p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">Source</p>
                  <EpEditRow ep={step.source} sources={sources} onChange={p => updSource(si, p)} />
                </div>
                <div>
                  <p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
                    Target{step.targets.length > 1 ? 's' : ''}
                  </p>
                  <div className="flex flex-col gap-2">
                    {step.targets.map((tgt, ti) => (
                      <div key={ti} className="flex items-center gap-2">
                        <EpEditRow ep={tgt} sources={sources} onChange={p => updTarget(si, ti, p)} />
                        {step.targets.length > 1 && (
                          <Button variant="ghost" size="icon" className="h-7 w-7 shrink-0"
                            onClick={() => removeTarget(si, ti)}>
                            <X className="h-3 w-3 text-muted-foreground" />
                          </Button>
                        )}
                      </div>
                    ))}
                    <button onClick={() => addTarget(si)}
                      className="text-xs text-blue-400 hover:text-blue-300 text-left mt-0.5">
                      + Add target
                    </button>
                  </div>
                </div>
                {si > 0 && (
                  <div>
                    <p className="text-[10px] uppercase tracking-widest text-muted-foreground mb-1.5">
                      Run after (step IDs, comma-separated)
                    </p>
                    <Input className="h-7 text-xs font-mono"
                      value={step.after?.join(', ') ?? ''}
                      placeholder="Leave blank to run in parallel"
                      onChange={e => updStep(si, { after: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })}
                    />
                  </div>
                )}
              </CardContent>
            </Card>
          ))}
          <button onClick={() => mutSteps(ss => [...ss, blankStep(firstSrcId)])}
            className="w-full rounded-xl border border-dashed border-border py-3 text-sm text-muted-foreground hover:border-blue-700 hover:text-blue-400 transition-colors">
            + Add Step
          </button>
        </div>
      </div>
    </div>
  )
}

// ── EpEditRow ──────────────────────────────────────────────────────────────

function EpEditRow({ ep, sources, onChange }:
  { ep: Endpoint; sources: Source[]; onChange: (p: Partial<Endpoint>) => void }) {
  return (
    <div className="flex gap-2 items-center">
      <Select value={ep.source_id} onValueChange={v => onChange({ source_id: v })}>
        <SelectTrigger className="w-36 h-7 text-xs shrink-0"><SelectValue placeholder="Source" /></SelectTrigger>
        <SelectContent>
          {sources.map(s => <SelectItem key={s.id} value={s.id}>{s.name}</SelectItem>)}
        </SelectContent>
      </Select>
      <Input className="h-7 text-xs w-24 shrink-0" value={ep.resource}
        onChange={e => onChange({ resource: e.target.value })} placeholder="resource" />
      <Input className="h-7 text-xs flex-1" value={ep.query ?? ''}
        onChange={e => onChange({ query: e.target.value })} placeholder="query (optional)" />
    </div>
  )
}
