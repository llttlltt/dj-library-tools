import { FlaskConical } from "lucide-react";
import { useEffect } from "react";
import type { QueryTesterOpts } from "@/App";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "@/components/ui/select";
import { cn } from "@/lib/utils";
import {
	filterWritableResources,
	findConnectionProvider,
} from "@/store/selection";
import type { Connection, Endpoint, ProviderInfo } from "@/types";

interface EndpointEditorProps {
	endpoint: Endpoint;
	connections: Connection[];
	providers: ProviderInfo[];
	isTarget?: boolean;
	onChange: (patch: Partial<Endpoint>) => void;
	onOpenQueryTester?: (opts?: QueryTesterOpts) => void;
	layout?: "row" | "grid";
}

export function EndpointEditor({
	endpoint,
	connections,
	providers,
	isTarget = false,
	onChange,
	onOpenQueryTester,
	layout = "row",
}: EndpointEditorProps) {
	const provider = findConnectionProvider(
		endpoint.connection_id,
		connections,
		providers,
	);

	// If this is a target, only show connections whose provider has at least one writable resource.
	const filteredConnections = connections.filter((c) => {
		if (!isTarget) return true;
		const p = providers.find((prov) => prov.name === c.provider);
		return p?.resources.some((r) => r.can_write) ?? true;
	});

	// For targets, only show resources that can be written to.
	const availableResources = filterWritableResources(provider, isTarget);

	// Automatically fix blank or invalid resource selections
	useEffect(() => {
		if (availableResources.length > 0) {
			const isValid = availableResources.some(
				(r) => r.name === endpoint.resource,
			);
			if (!isValid) {
				onChange({ resource: availableResources[0].name });
			}
		}
	}, [endpoint.resource, availableResources, onChange]);

	const currentRes = provider?.resources.find(
		(r) => r.name === endpoint.resource,
	);
	const supportsQuery = currentRes?.supports_query ?? true;
	const isInvalidTarget = isTarget && currentRes && !currentRes.can_write;

	const connectionSelect = (
		<div
			className={cn(
				"flex flex-col gap-1.5 min-w-0",
				layout === "row" ? "shrink-0" : "flex-1",
			)}
		>
			{layout === "grid" && (
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Connection
				</span>
			)}
			<Select
				value={endpoint.connection_id}
				onValueChange={(v) => {
					const newConnection = connections.find((c) => c.id === v);
					const newProv = providers.find(
						(p) => p.name === newConnection?.provider,
					);
					const newResList = filterWritableResources(newProv, isTarget);
					const nextRes = newResList[0]?.name ?? endpoint.resource;
					onChange({ connection_id: v, resource: nextRes });
				}}
			>
				<SelectTrigger
					className={cn(
						"text-sm bg-background/50",
						layout === "row" ? "w-40 h-9 shrink-0" : "h-8.5",
					)}
				>
					<SelectValue placeholder="Select a connection…" />
				</SelectTrigger>
				<SelectContent>
					{[...filteredConnections]
						.sort((a, b) => a.name.localeCompare(b.name))
						.map((c) => (
							<SelectItem key={c.id} value={c.id}>
								{c.name}
							</SelectItem>
						))}
				</SelectContent>
			</Select>
		</div>
	);

	const resourceSelect = (
		<div className="flex flex-col gap-1.5 shrink-0">
			{layout === "grid" && (
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Resource
				</span>
			)}
			<Select
				value={endpoint.resource}
				onValueChange={(v) => onChange({ resource: v })}
			>
				<SelectTrigger
					className={cn(
						"text-sm bg-background/50 transition-colors",
						layout === "row" ? "w-24 h-9" : "h-8.5",
						isInvalidTarget &&
							"border-destructive/60 bg-destructive/5 text-destructive",
					)}
				>
					<SelectValue placeholder="resource" />
				</SelectTrigger>
				<SelectContent>
					{availableResources.map((r) => (
						<SelectItem key={r.name} value={r.name}>
							{r.name}
						</SelectItem>
					))}
				</SelectContent>
			</Select>
			{isInvalidTarget && layout === "grid" && (
				<span className="text-[10px] text-destructive/80 font-medium animate-in fade-in slide-in-from-top-1">
					Read-only resource.
				</span>
			)}
		</div>
	);

	const queryInput = (
		<div className="flex flex-col gap-1.5 flex-1 min-w-0">
			{layout === "grid" && (
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Query Expression
				</span>
			)}
			<div className="relative">
				<Input
					className={cn(
						"text-sm w-full bg-background/50 transition-opacity font-mono",
						layout === "row" ? "h-9" : "h-8.5",
						supportsQuery ? "opacity-100" : "opacity-10 pointer-events-none",
					)}
					value={endpoint.query ?? ""}
					onChange={(e) => onChange({ query: e.target.value })}
					placeholder={
						supportsQuery
							? layout === "row"
								? "query (optional)"
								: "beatgrids-count:1 && bpm:>120"
							: "no query supported"
					}
					disabled={!supportsQuery}
				/>
			</div>
		</div>
	);

	if (layout === "grid") {
		return (
			<div className="space-y-4">
				<div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
					{connectionSelect}
					{resourceSelect}
				</div>
				{queryInput}
			</div>
		);
	}

	return (
		<div className="flex flex-nowrap gap-2.5 items-center w-full min-w-0">
			{connectionSelect}
			{resourceSelect}
			{queryInput}

			{onOpenQueryTester && supportsQuery ? (
				<Button
					type="button"
					variant="ghost"
					size="icon"
					className="h-8.5 w-8.5 shrink-0 hover:bg-secondary"
					title="Test query"
					onClick={() =>
						onOpenQueryTester({
							connectionID: endpoint.connection_id,
							resource: endpoint.resource,
							query: endpoint.query ?? "",
							isTarget,
							onApply: (q) => onChange({ query: q }),
						})
					}
				>
					<FlaskConical className="h-4 w-4 text-muted-foreground" />
				</Button>
			) : (
				<div className="w-8.5 shrink-0" />
			)}
		</div>
	);
}
