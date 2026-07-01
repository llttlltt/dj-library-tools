import { useAtom } from "@effect-atom/atom-react";
import { Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { SourceCard } from "@/components/source/SourceCard";
import { type Provider, SourceForm } from "@/components/source/SourceForm";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { runPromise } from "@/lib/runtime";
import {
	addSource,
	loadSources,
	removeSource,
	sourcesAtom,
	sourcesErrorAtom,
	updateSource,
} from "@/store/sources";
import type { Source } from "@/types";

export default function SourcesView() {
	const [sources] = useAtom(sourcesAtom);
	const [error] = useAtom(sourcesErrorAtom);
	const [formOpen, setFormOpen] = useState(false);
	const [editTarget, setEditTarget] = useState<Source | null>(null);
	const [deleteTarget, setDeleteTarget] = useState<Source | null>(null);

	useEffect(() => {
		runPromise(loadSources);
	}, []);

	async function handleSubmit(
		name: string,
		provider: Provider,
		cfg: Record<string, string>,
		id?: string,
	) {
		const src: Source = { id: id ?? "", name, provider, config: cfg };
		if (id) {
			await runPromise(updateSource(src));
		} else {
			await runPromise(addSource(name, provider, cfg));
		}
	}

	async function confirmDelete(s: Source) {
		setDeleteTarget(null);
		await runPromise(removeSource(s.id));
	}

	return (
		<div className="flex flex-col h-full">
			{/* Sticky Top Header Nav */}
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

			{/* Scrollable Main Content Box */}
			<div className="flex-1 overflow-auto p-6">
				{error && <p className="text-sm text-destructive mb-4">{error}</p>}

				{sources.length === 0 ? (
					<p className="text-sm text-muted-foreground italic">
						No sources configured. Click "Add Source" to get started.
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
