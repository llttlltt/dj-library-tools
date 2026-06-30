import { FolderOpen } from "lucide-react";
import type React from "react";
import { useEffect, useState } from "react";
import { PlexAuthModal } from "@/components/source/PlexAuthModal";
import { Button } from "@/components/ui/button";
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
import { OpenFileDialog } from "../../../wailsjs/go/gui/App";

export type Provider = "rb" | "m3u" | "plex";

interface SourceFormProps {
	open: boolean;
	initial?: Source | null;
	onClose: () => void;
	onSubmit: (
		name: string,
		provider: Provider,
		cfg: Record<string, string>,
		id?: string,
	) => Promise<void>;
}

export function SourceForm({
	open,
	initial,
	onClose,
	onSubmit,
}: SourceFormProps) {
	const isEdit = !!initial;

	const [name, setName] = useState(initial?.name ?? "");
	const [provider, setProvider] = useState<Provider>(
		(initial?.provider as Provider) ?? "rb",
	);
	const [filePath, setFilePath] = useState(initial?.config?.file_path ?? "");
	const [host, setHost] = useState(initial?.config?.host ?? "");
	const [port, setPort] = useState(initial?.config?.port ?? "");
	const [token, setToken] = useState(initial?.config?.token ?? "");
	const [error, setError] = useState("");
	const [saving, setSaving] = useState(false);
	const [plexAuth, setPlexAuth] = useState(false);

	// Reset all controlled fields whenever the dialog opens or the target source changes.
	useEffect(() => {
		if (open) {
			setName(initial?.name ?? "");
			setProvider((initial?.provider as Provider) ?? "rb");
			setFilePath(initial?.config?.file_path ?? "");
			setHost(initial?.config?.host ?? "");
			setPort(initial?.config?.port ?? "");
			setToken(initial?.config?.token ?? "");
			setError("");
		}
	}, [open, initial]);

	const resetToInitial = () => {
		setName(initial?.name ?? "");
		setProvider((initial?.provider as Provider) ?? "rb");
		setFilePath(initial?.config?.file_path ?? "");
		setHost(initial?.config?.host ?? "");
		setPort(initial?.config?.port ?? "");
		setToken(initial?.config?.token ?? "");
		setError("");
	};

	async function pickFile() {
		try {
			const dir = filePath.substring(0, filePath.lastIndexOf("/")) || "";
			const path = await OpenFileDialog(dir);
			if (path) setFilePath(path);
		} catch (e) {
			setError(String(e));
		}
	}

	async function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		setError("");
		setSaving(true);
		const cfg: Record<string, string> =
			provider === "plex" ? { host, port, token } : { file_path: filePath };
		try {
			await onSubmit(name, provider, cfg, initial?.id);
			onClose();
		} catch (err) {
			setError(String(err));
		}
		setSaving(false);
	}

	return (
		<>
			<Dialog
				open={open}
				onOpenChange={(o) => {
					if (!o) {
						onClose();
						resetToInitial();
					}
				}}
			>
				<DialogContent>
					<DialogHeader>
						<DialogTitle>{isEdit ? "Edit Source" : "Add Source"}</DialogTitle>
					</DialogHeader>
					<form onSubmit={handleSubmit} className="flex flex-col gap-4 mt-4">
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
								<div className="flex gap-2">
									<Input
										id="src-filepath"
										value={filePath}
										onChange={(e) => setFilePath(e.target.value)}
										placeholder="/path/to/library.xml"
										required
										className="flex-1"
									/>
									<Button
										type="button"
										variant="outline"
										size="icon"
										onClick={pickFile}
										title="Browse for file"
									>
										<FolderOpen className="h-4 w-4" />
									</Button>
								</div>
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
										{token ? (
											<div className="flex gap-2 items-center">
												<Input
													id="src-token"
													value={token}
													onChange={(e) => setToken(e.target.value)}
													placeholder="plex-token"
												/>
											</div>
										) : (
											<Button
												type="button"
												variant="outline"
												size="sm"
												onClick={() => setPlexAuth(true)}
											>
												Authenticate with Plex…
											</Button>
										)}
									</div>
								</div>
								{token && (
									<button
										type="button"
										className="text-xs text-blue-400 hover:text-blue-300 text-left -mt-1"
										onClick={() => setPlexAuth(true)}
									>
										Re-authenticate with Plex
									</button>
								)}
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
								{saving ? "Saving…" : isEdit ? "Update Source" : "Save Source"}
							</Button>
						</div>
					</form>
				</DialogContent>
			</Dialog>

			{/* Plex PIN auth modal — outside the main Dialog to avoid nesting */}
			<PlexAuthModal
				open={plexAuth}
				onToken={(t) => {
					setToken(t);
					setPlexAuth(false);
				}}
				onCancel={() => setPlexAuth(false)}
			/>
		</>
	);
}
