import { Plus } from "lucide-react";
import { useCallback, useEffect, useState } from "react";
import { SourceCard } from "@/components/source/SourceCard";
import { type Provider, SourceForm } from "@/components/source/SourceForm";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import type { Source } from "@/types";
import {
	CreateSource,
	DeleteSource,
	ListSources,
	UpdateSource,
} from "../../wailsjs/go/gui/App";

export default function SourcesView() {
	const [sources, setSources] = useState<Source[]>([]);
	const [error, setError] = useState("");
	const [formOpen, setFormOpen] = useState(false);
	const [editTarget, setEditTarget] = useState<Source | null>(null);
	const [deleteTarget, setDeleteTarget] = useState<Source | null>(null);

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

	async function handleSubmit(
		name: string,
		provider: Provider,
		cfg: Record<string, string>,
		id?: string,
	) {
		const src: Source = { id: id ?? "", name, provider, config: cfg };
		if (id) {
			await (UpdateSource as (s: unknown) => Promise<void>)(src);
		} else {
			await CreateSource(name, provider, cfg);
		}
		await load();
	}

	async function confirmDelete(s: Source) {
		setDeleteTarget(null);
		try {
			await DeleteSource(s.id);
			await load();
		} catch (e) {
			setError(String(e));
		}
	}

	return (
		<div className="flex flex-col h-full">
			<div className="h-14 flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
				<span className="text-sm font-semibold">Sources</span>
				<div className="flex-1" />
				<Button
					type="button"
					size="sm"
					onClick={() => {
						setEditTarget(null);
						setFormOpen(true);
					}}
				>
					<Plus className="h-4 w-4 mr-1.5" /> Add Source
				</Button>
			</div>
			<div className="flex-1 overflow-auto p-6">
				{error && <p className="text-sm text-destructive mb-4">{error}</p>}

				{sources.length === 0 ? (
					<p className="text-sm text-muted-foreground italic">
						No sources configured. Add one to get started.
					</p>
				) : (
					<div className="flex flex-col gap-3">
						{[...sources]
							.sort((a, b) => a.name.localeCompare(b.name))
							.map((s) => (
								<SourceCard
									key={s.id}
									source={s}
									onEdit={(src) => {
										setEditTarget(src);
										setFormOpen(true);
									}}
									onDelete={(src) => setDeleteTarget(src)}
								/>
							))}
					</div>
				)}

				<SourceForm
					open={formOpen}
					initial={editTarget}
					onClose={() => {
						setFormOpen(false);
						setEditTarget(null);
					}}
					onSubmit={handleSubmit}
				/>

				{deleteTarget && (
					<ConfirmDialog
						open={true}
						title={`Delete "${deleteTarget.name}"?`}
						description="This source will be permanently removed."
						confirmLabel="Delete"
						destructive
						onConfirm={() => confirmDelete(deleteTarget)}
						onCancel={() => setDeleteTarget(null)}
					/>
				)}
			</div>
		</div>
	);
}
