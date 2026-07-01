import { useAtom } from "@effect-atom/atom-react";
import { Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { ConnectionCard } from "@/components/connection/ConnectionCard";
import {
	ConnectionForm,
	type Provider,
} from "@/components/connection/ConnectionForm";
import { Button } from "@/components/ui/button";
import { ConfirmDialog } from "@/components/ui/confirm-dialog";
import { runPromise } from "@/lib/runtime";
import {
	addConnection,
	connectionsAtom,
	connectionsErrorAtom,
	loadConnections,
	removeConnection,
	updateConnection,
} from "@/store/connections";
import type { Connection } from "@/types";

export default function ConnectionsView() {
	const [connections] = useAtom(connectionsAtom);
	const [error] = useAtom(connectionsErrorAtom);
	const [formOpen, setFormOpen] = useState(false);
	const [editTarget, setEditTarget] = useState<Connection | null>(null);
	const [deleteTarget, setDeleteTarget] = useState<Connection | null>(null);

	useEffect(() => {
		runPromise(loadConnections);
	}, []);

	async function handleSubmit(
		name: string,
		provider: Provider,
		cfg: Record<string, string>,
		id?: string,
	) {
		const conn: Connection = { id: id ?? "", name, provider, config: cfg };
		if (id) {
			await runPromise(updateConnection(conn));
		} else {
			await runPromise(addConnection(name, provider, cfg));
		}
	}

	async function confirmDelete(c: Connection) {
		setDeleteTarget(null);
		await runPromise(removeConnection(c.id));
	}

	return (
		<div className="flex flex-col h-full">
			{/* Sticky Top Header Nav */}
			<div className="h-14 flex items-center gap-2 px-6 py-3 border-b border-border bg-[hsl(240_10%_4%)] sticky top-0 z-10">
				<span className="text-sm font-semibold">Connections</span>
				<div className="flex-1" />
				<Button
					type="button"
					size="sm"
					onClick={() => {
						setEditTarget(null);
						setFormOpen(true);
					}}
				>
					<Plus className="h-4 w-4 mr-1.5" /> Add Connection
				</Button>
			</div>

			{/* Scrollable Main Content Box */}
			<div className="flex-1 overflow-auto p-6">
				{error && <p className="text-sm text-destructive mb-4">{error}</p>}

				{connections.length === 0 ? (
					<p className="text-sm text-muted-foreground italic">
						No connections configured. Click "Add Connection" to get started.
					</p>
				) : (
					<div className="flex flex-col gap-3">
						{[...connections]
							.sort((a, b) => a.name.localeCompare(b.name))
							.map((c) => (
								<ConnectionCard
									key={c.id}
									connection={c}
									onEdit={(conn) => {
										setEditTarget(conn);
										setFormOpen(true);
									}}
									onDelete={(conn) => setDeleteTarget(conn)}
								/>
							))}
					</div>
				)}

				<ConnectionForm
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
						description="This connection will be permanently removed."
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
