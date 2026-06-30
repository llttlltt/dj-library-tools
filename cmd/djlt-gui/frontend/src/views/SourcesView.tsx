import { Plus, Trash2 } from "lucide-react";
import type React from "react";
import { useCallback, useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import type { Source } from "@/types";
import {
	CreateSource,
	DeleteSource,
	ListSources,
} from "../../wailsjs/go/gui/App";

type Provider = "rb" | "m3u" | "plex";
const PROVIDER_LABELS: Record<string, string> = {
	rb: "Rekordbox",
	m3u: "M3U",
	plex: "Plex",
};

export default function SourcesView() {
	const [sources, setSources] = useState<Source[]>([]);
	const [open, setOpen] = useState(false);
	const [error, setError] = useState("");
	const [saving, setSaving] = useState(false);

	const [name, setName] = useState("");
	const [provider, setProvider] = useState<Provider>("rb");
	const [filePath, setFilePath] = useState("");
	const [host, setHost] = useState("");
	const [port, setPort] = useState("");
	const [token, setToken] = useState("");

	const load = useCallback(async () => {
		try {
			setSources(((await ListSources()) as unknown as Source[]) ?? []);
		} catch (e) {
			setError(String(e));
		}
	}, []);

	useEffect(() => {
		load();
	}, [load]);

	function resetForm() {
		setName("");
		setProvider("rb");
		setFilePath("");
		setHost("");
		setPort("");
		setToken("");
		setError("");
	}

	async function handleCreate(e: React.FormEvent) {
		e.preventDefault();
		setError("");
		setSaving(true);
		const cfg: Record<string, string> =
			provider === "plex" ? { host, port, token } : { file_path: filePath };
		try {
			await CreateSource(name, provider, cfg);
			setOpen(false);
			resetForm();
			await load();
		} catch (e) {
			setError(String(e));
		}
		setSaving(false);
	}

	async function handleDelete(id: string, label: string) {
		if (!confirm(`Delete source "${label}"?`)) return;
		try {
			await DeleteSource(id);
			await load();
		} catch (e) {
			setError(String(e));
		}
	}

	return (
		<div className="p-6 max-w-2xl">
			<div className="flex items-center justify-between mb-6">
				<div>
					<h1 className="text-lg font-semibold">Sources</h1>
					<p className="text-sm text-muted-foreground mt-0.5">
						Provider connections for your DJ libraries
					</p>
				</div>
				<Button
					type="button"
					size="sm"
					onClick={() => {
						resetForm();
						setOpen(true);
					}}
				>
					<Plus className="h-4 w-4 mr-1.5" /> Add Source
				</Button>
			</div>

			{error && <p className="text-sm text-destructive mb-4">{error}</p>}

			{sources.length === 0 ? (
				<p className="text-sm text-muted-foreground italic">
					No sources configured. Add one to get started.
				</p>
			) : (
				<div className="flex flex-col gap-3">
					{sources.map((s) => (
						<Card key={s.id}>
							<CardHeader className="flex-row items-center justify-between py-3 px-4 gap-0">
								<div className="flex items-center gap-3">
									<CardTitle className="text-sm">{s.name}</CardTitle>
									<Badge
										variant={
											s.provider === "rb"
												? "sync"
												: s.provider === "plex"
													? "fix"
													: "edit"
										}
									>
										{PROVIDER_LABELS[s.provider] ?? s.provider}
									</Badge>
								</div>
								<Button
									type="button"
									variant="ghost"
									size="icon"
									onClick={() => handleDelete(s.id, s.name)}
								>
									<Trash2 className="h-3.5 w-3.5 text-muted-foreground hover:text-destructive" />
								</Button>
							</CardHeader>
							<CardContent className="py-0 pb-3 px-4">
								<p className="text-xs text-muted-foreground font-mono">
									{s.id}
								</p>
								{s.config?.file_path && (
									<p className="text-xs text-muted-foreground mt-0.5 truncate">
										{s.config.file_path}
									</p>
								)}
							</CardContent>
						</Card>
					))}
				</div>
			)}

			<Dialog
				open={open}
				onOpenChange={(o) => {
					setOpen(o);
					if (!o) resetForm();
				}}
			>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>Add Source</DialogTitle>
					</DialogHeader>
					<form onSubmit={handleCreate} className="flex flex-col gap-4 mt-4">
						<div className="flex flex-col gap-1.5">
							<label
								htmlFor="src-name"
								className="text-xs text-muted-foreground uppercase tracking-wide"
							>
								Name
							</label>
							<Input
								id="src-name"
								value={name}
								onChange={(e) => setName(e.target.value)}
								placeholder="Main Library"
								required
							/>
						</div>
						<div className="flex flex-col gap-1.5">
							<label
								htmlFor="src-provider"
								className="text-xs text-muted-foreground uppercase tracking-wide"
							>
								Provider
							</label>
							<Select
								value={provider}
								onValueChange={(v) => setProvider(v as Provider)}
							>
								<SelectTrigger id="src-provider">
									<SelectValue />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="rb">Rekordbox</SelectItem>
									<SelectItem value="m3u">M3U</SelectItem>
									<SelectItem value="plex">Plex</SelectItem>
								</SelectContent>
							</Select>
						</div>
						{(provider === "rb" || provider === "m3u") && (
							<div className="flex flex-col gap-1.5">
								<label
									htmlFor="src-filepath"
									className="text-xs text-muted-foreground uppercase tracking-wide"
								>
									File Path
								</label>
								<Input
									id="src-filepath"
									value={filePath}
									onChange={(e) => setFilePath(e.target.value)}
									placeholder="/path/to/library.xml"
									required
								/>
							</div>
						)}
						{provider === "plex" && (
							<>
								<div className="flex flex-col gap-1.5">
									<label
										htmlFor="src-host"
										className="text-xs text-muted-foreground uppercase tracking-wide"
									>
										Host
									</label>
									<Input
										id="src-host"
										value={host}
										onChange={(e) => setHost(e.target.value)}
										placeholder="localhost"
									/>
								</div>
								<div className="flex gap-3">
									<div className="flex flex-col gap-1.5 flex-1">
										<label
											htmlFor="src-port"
											className="text-xs text-muted-foreground uppercase tracking-wide"
										>
											Port
										</label>
										<Input
											id="src-port"
											value={port}
											onChange={(e) => setPort(e.target.value)}
											placeholder="32400"
										/>
									</div>
									<div className="flex flex-col gap-1.5 flex-1">
										<label
											htmlFor="src-token"
											className="text-xs text-muted-foreground uppercase tracking-wide"
										>
											Token
										</label>
										<Input
											id="src-token"
											value={token}
											onChange={(e) => setToken(e.target.value)}
											placeholder="plex-token"
										/>
									</div>
								</div>
							</>
						)}
						{error && <p className="text-sm text-destructive">{error}</p>}
						<div className="flex justify-end gap-2 mt-1">
							<DialogClose asChild>
								<Button type="button" variant="outline" size="sm">
									Cancel
								</Button>
							</DialogClose>
							<Button type="submit" size="sm" disabled={saving}>
								{saving ? "Saving…" : "Save Source"}
							</Button>
						</div>
					</form>
				</DialogContent>
			</Dialog>
		</div>
	);
}
