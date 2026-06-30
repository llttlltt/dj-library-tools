import React, { useEffect, useState } from 'react'
import { Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogClose } from '@/components/ui/dialog'
import { ListSources, CreateSource, DeleteSource } from '../../wailsjs/go/gui/App'
import type { Source } from '@/types'

type Provider = 'rb' | 'm3u' | 'plex'
const PROVIDER_LABELS: Record<string, string> = { rb: 'Rekordbox', m3u: 'M3U', plex: 'Plex' }

export default function SourcesView() {
  const [sources, setSources] = useState<Source[]>([])
  const [open,    setOpen]    = useState(false)
  const [error,   setError]   = useState('')
  const [saving,  setSaving]  = useState(false)

  const [name,     setName]     = useState('')
  const [provider, setProvider] = useState<Provider>('rb')
  const [filePath, setFilePath] = useState('')
  const [host,     setHost]     = useState('')
  const [port,     setPort]     = useState('')
  const [token,    setToken]    = useState('')

  async function load() {
    try { setSources((await ListSources() as unknown as Source[]) ?? []) }
    catch (e) { setError(String(e)) }
  }

  useEffect(() => { load() }, [])

  function resetForm() {
    setName(''); setProvider('rb'); setFilePath('')
    setHost(''); setPort(''); setToken(''); setError('')
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault(); setError(''); setSaving(true)
    const cfg: Record<string, string> =
      provider === 'plex'
        ? { host, port, token }
        : { file_path: filePath }
    try { await CreateSource(name, provider, cfg); setOpen(false); resetForm(); await load() }
    catch (e) { setError(String(e)) }
    setSaving(false)
  }

  async function handleDelete(id: string, label: string) {
    if (!confirm(`Delete source "${label}"?`)) return
    try { await DeleteSource(id); await load() }
    catch (e) { setError(String(e)) }
  }

  return (
    <div className="p-6 max-w-2xl">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-lg font-semibold">Sources</h1>
          <p className="text-sm text-muted-foreground mt-0.5">Provider connections for your DJ libraries</p>
        </div>
        <Button size="sm" onClick={() => { resetForm(); setOpen(true) }}>
          <Plus className="h-4 w-4 mr-1.5" /> Add Source
        </Button>
      </div>

      {error && <p className="text-sm text-destructive mb-4">{error}</p>}

      {sources.length === 0
        ? <p className="text-sm text-muted-foreground italic">No sources configured. Add one to get started.</p>
        : <div className="flex flex-col gap-3">
            {sources.map(s => (
              <Card key={s.id}>
                <CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
                  <div className="flex items-center gap-3">
                    <CardTitle className="text-sm">{s.name}</CardTitle>
                    <Badge variant={s.provider === 'rb' ? 'sync' : s.provider === 'plex' ? 'fix' : 'edit'}>
                      {PROVIDER_LABELS[s.provider] ?? s.provider}
                    </Badge>
                  </div>
                  <Button variant="ghost" size="icon" onClick={() => handleDelete(s.id, s.name)}>
                    <Trash2 className="h-3.5 w-3.5 text-muted-foreground hover:text-destructive" />
                  </Button>
                </CardHeader>
                <CardContent className="py-0 pb-3 px-4">
                  <p className="text-xs text-muted-foreground font-mono">{s.id}</p>
                  {s.config?.file_path && (
                    <p className="text-xs text-muted-foreground mt-0.5 truncate">{s.config.file_path}</p>
                  )}
                </CardContent>
              </Card>
            ))}
          </div>
      }

      <Dialog open={open} onOpenChange={o => { setOpen(o); if (!o) resetForm() }}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Add Source</DialogTitle>
          </DialogHeader>
          <form onSubmit={handleCreate} className="flex flex-col gap-4 mt-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-xs text-muted-foreground uppercase tracking-wide">Name</label>
              <Input value={name} onChange={e => setName(e.target.value)} placeholder="Main Library" required />
            </div>
            <div className="flex flex-col gap-1.5">
              <label className="text-xs text-muted-foreground uppercase tracking-wide">Provider</label>
              <Select value={provider} onValueChange={v => setProvider(v as Provider)}>
                <SelectTrigger><SelectValue /></SelectTrigger>
                <SelectContent>
                  <SelectItem value="rb">Rekordbox</SelectItem>
                  <SelectItem value="m3u">M3U</SelectItem>
                  <SelectItem value="plex">Plex</SelectItem>
                </SelectContent>
              </Select>
            </div>
            {(provider === 'rb' || provider === 'm3u') && (
              <div className="flex flex-col gap-1.5">
                <label className="text-xs text-muted-foreground uppercase tracking-wide">File Path</label>
                <Input value={filePath} onChange={e => setFilePath(e.target.value)} placeholder="/path/to/library.xml" required />
              </div>
            )}
            {provider === 'plex' && (
              <>
                <div className="flex flex-col gap-1.5">
                  <label className="text-xs text-muted-foreground uppercase tracking-wide">Host</label>
                  <Input value={host} onChange={e => setHost(e.target.value)} placeholder="localhost" />
                </div>
                <div className="flex gap-3">
                  <div className="flex flex-col gap-1.5 flex-1">
                    <label className="text-xs text-muted-foreground uppercase tracking-wide">Port</label>
                    <Input value={port} onChange={e => setPort(e.target.value)} placeholder="32400" />
                  </div>
                  <div className="flex flex-col gap-1.5 flex-1">
                    <label className="text-xs text-muted-foreground uppercase tracking-wide">Token</label>
                    <Input value={token} onChange={e => setToken(e.target.value)} placeholder="plex-token" />
                  </div>
                </div>
              </>
            )}
            {error && <p className="text-sm text-destructive">{error}</p>}
            <div className="flex justify-end gap-2 mt-1">
              <DialogClose asChild>
                <Button type="button" variant="outline" size="sm">Cancel</Button>
              </DialogClose>
              <Button type="submit" size="sm" disabled={saving}>{saving ? 'Saving…' : 'Save Source'}</Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>
    </div>
  )
}
