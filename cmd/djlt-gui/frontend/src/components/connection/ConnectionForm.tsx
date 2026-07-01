import { useAtom } from "@effect-atom/atom-react";
import { Loader2, Plus, X } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { runPromise } from "@/lib/runtime";
import { AppService } from "@/services";
import { providersAtom } from "@/store/providers";
import type { Connection } from "@/types";

export type Provider = "rb" | "plex" | "m3u" | "m3u8";

interface ConnectionFormProps {
	open: boolean;
	initial?: Connection | null;
	onClose: () => void;
	onSubmit: (
		name: string,
		provider: Provider,
		config: Record<string, string>,
		id?: string,
	) => Promise<void>;
}

export function ConnectionForm({
	open,
	initial,
	onClose,
	onSubmit,
}: ConnectionFormProps) {
	const [providers] = useAtom(providersAtom);
	const [name, setName] = useState(initial?.name ?? "");
	const [provider, setProvider] = useState<Provider>(
		(initial?.provider as Provider) ?? "rb",
	);
	const [config, setConfig] = useState<Record<string, string>>(
		initial?.config ?? {},
	);
	const [saving, setSaving] = useState(false);

	const isEdit = !!initial;

	useEffect(() => {
		if (initial) {
			setName(initial.name);
			setProvider((initial.provider as Provider) ?? "rb");
			setConfig(initial.config);
		} else {
			setName("");
			setProvider("rb");
			setConfig({});
		}
	}, [initial]);

	async function handlePickFile(key: string) {
		const app = await runPromise(AppService);
		const path = await runPromise(app.openFileDialog(""));
		if (path) {
			setConfig((prev) => ({ ...prev, [key]: path as string }));
		}
	}

	async function handlePlexAuth() {
		const app = await runPromise(AppService);
		const auth = await runPromise(app.initPlexAuth());
		window.open(auth.url, "_blank");

		// Poll for token
		const interval = setInterval(async () => {
			const token = await runPromise(app.checkPlexAuth(auth.pin_id as number));
			if (token) {
				setConfig((prev) => ({ ...prev, token: token as string }));
				clearInterval(interval);
			}
		}, 2000);
	}

	async function handleInternalSubmit(e: React.FormEvent) {
		e.preventDefault();
		setSaving(true);
		try {
			await onSubmit(name, provider, config, initial?.id);
			onClose();
		} finally {
			setSaving(false);
		}
	}

	const selectedProviderInfo = providers.find((p) => p.name === provider);

	return (
		<Dialog open={open} onOpenChange={(o) => !o && onClose()}>
			<DialogContent className="sm:max-w-[425px]">
				<form onSubmit={handleInternalSubmit}>
					<DialogHeader>
						<DialogTitle>
							{isEdit ? "Edit Connection" : "Add Connection"}
						</DialogTitle>
						<DialogDescription>
							Configure how the application connects to your DJ library.
						</DialogDescription>
					</DialogHeader>

					<div className="grid gap-5 py-6">
						<div className="grid gap-2">
							<Label htmlFor="name">Display Name</Label>
							<Input
								id="name"
								value={name}
								onChange={(e) => setName(e.target.value)}
								placeholder="e.g. My Main Library"
								required
							/>
						</div>

						<div className="grid gap-2">
							<Label htmlFor="provider">Provider</Label>
							<Select
								value={provider}
								onValueChange={(v) => setProvider(v as Provider)}
								disabled={isEdit}
							>
								<SelectTrigger>
									<SelectValue placeholder="Select provider" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem value="rb">Rekordbox XML</SelectItem>
									<SelectItem value="plex">Plex</SelectItem>
									<SelectItem value="m3u">M3U Playlist</SelectItem>
									<SelectItem value="m3u8">M3U8 (UTF-8) Playlist</SelectItem>
								</SelectContent>
							</Select>
						</div>

						{/* Dynamic config based on provider */}
						{selectedProviderInfo?.capabilities.IsFileBased && (
							<div className="grid gap-2">
								<Label>Library File</Label>
								<div className="flex gap-2">
									<Input
										value={config.file_path || ""}
										readOnly
										placeholder="Select a file…"
										className="bg-muted/50"
										required
									/>
									<Button
										type="button"
										variant="outline"
										onClick={() => handlePickFile("file_path")}
									>
										Browse
									</Button>
								</div>
							</div>
						)}

						{provider === "plex" && (
							<div className="space-y-4">
								<div className="grid gap-2">
									<Label htmlFor="host">Server Host</Label>
									<Input
										id="host"
										value={config.host || ""}
										onChange={(e) =>
											setConfig((prev) => ({ ...prev, host: e.target.value }))
										}
										placeholder="localhost"
										required
									/>
								</div>
								<div className="grid gap-2">
									<Label htmlFor="port">Port</Label>
									<Input
										id="port"
										value={config.port || "32400"}
										onChange={(e) =>
											setConfig((prev) => ({ ...prev, port: e.target.value }))
										}
										placeholder="32400"
										required
									/>
								</div>
								<div className="grid gap-2">
									<Label>Authentication</Label>
									{config.token ? (
										<div className="flex items-center gap-2 text-sm text-green-500 bg-green-500/10 p-2 rounded-md border border-green-500/20">
											<Plus className="w-4 h-4 rotate-45" /> Authenticated
											<Button
												variant="ghost"
												size="icon"
												className="h-6 w-6 ml-auto"
												onClick={() =>
													setConfig((prev) => {
														const next = { ...prev };
														delete next.token;
														return next;
													})
												}
											>
												<X className="w-3 h-3" />
											</Button>
										</div>
									) : (
										<Button
											type="button"
											variant="secondary"
											onClick={handlePlexAuth}
											className="w-full"
										>
											Link Plex Account
										</Button>
									)}
								</div>
							</div>
						)}
					</div>

					<DialogFooter>
						<Button type="button" variant="ghost" onClick={onClose}>
							Cancel
						</Button>
						<Button type="submit" disabled={saving}>
							{saving ? (
								<Loader2 className="w-4 h-4 mr-2 animate-spin" />
							) : null}
							{saving
								? "Saving…"
								: isEdit
									? "Update Connection"
									: "Save Connection"}
						</Button>
					</DialogFooter>
				</form>
			</DialogContent>
		</Dialog>
	);
}
