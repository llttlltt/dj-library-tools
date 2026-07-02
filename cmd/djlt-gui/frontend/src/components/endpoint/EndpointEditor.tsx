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
import { filterWritableResources } from "@/store/selection";
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

// ── Sub-components ─────────────────────────────────────────────────────────

function AdHocResourceControl({
	resource,
	onChange,
	layout,
}: {
	resource: string;
	onChange: (v: string) => void;
	layout: "row" | "grid";
}) {
	return (
		<div className="flex flex-col gap-1.5 min-w-0 flex-[0.6]">
			{layout === "grid" && (
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					File Path
				</span>
			)}
			<Input
				className={cn(
					"text-xs bg-background/50 font-mono",
					layout === "row" ? "h-9 w-40 shrink-0" : "h-8.5",
				)}
				value={resource}
				onChange={(e) => onChange(e.target.value)}
				placeholder="/path/to/file.m3u"
			/>
		</div>
	);
}

function StandardResourceControl({
	resource,
	availableResources,
	isInvalidTarget,
	onChange,
	layout,
}: {
	resource: string;
	availableResources: { name: string }[];
	isInvalidTarget: boolean | undefined;
	onChange: (v: string) => void;
	layout: "row" | "grid";
}) {
	return (
		<div className="flex flex-col gap-1.5 shrink-0">
			{layout === "grid" && (
				<span className="text-[10px] font-semibold uppercase tracking-wider text-muted-foreground">
					Resource
				</span>
			)}
			<Select value={resource} onValueChange={onChange}>
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
}

// ── Main Component ──────────────────────────────────────────────────────────

export function EndpointEditor({
	endpoint,
	connections,
	providers,
	isTarget = false,
	onChange,
	onOpenQueryTester,
	layout = "row",
}: EndpointEditorProps) {
	// If this is a target, only show connections whose provider has at least one writable resource.
	const filteredConnections = connections.filter((c) => {
		if (!isTarget) return true;
		const p = providers.find((prov) => prov.name === c.provider);
		return p?.resources.some((r) => r.can_write) ?? true;
	});

	const adHoc: Connection[] = [
		{ id: "m3u", name: "AD-HOC M3U", provider: "m3u", config: {} },
		{ id: "m3u8", name: "AD-HOC M3U8", provider: "m3u8", config: {} },
	];

	const displayConnections = [...filteredConnections, ...adHoc];

	const selectedConnection = displayConnections.find(
		(c) => c.id === endpoint.connection_id,
	);
	const provider = providers.find(
		(p) => p.name === selectedConnection?.provider,
	);

	const isAdHoc =
		endpoint.connection_id === "m3u" || endpoint.connection_id === "m3u8";

	// For targets, only show resources that can be written to.
	const availableResources = filterWritableResources(provider, isTarget);

	// Automatically fix blank or invalid resource selections
	useEffect(() => {
		if (!isAdHoc && availableResources.length > 0) {
			const isValid = availableResources.some(
				(r) => r.name === endpoint.resource,
			);
			if (!isValid) {
				onChange({ resource: availableResources[0].name });
			}
		}
	}, [endpoint.resource, availableResources, onChange, isAdHoc]);

	const currentRes = provider?.resources.find(
		(r) => r.name === endpoint.resource,
	);
	const supportsQuery = isAdHoc || (currentRes?.supports_query ?? true);
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
					const newConnection = displayConnections.find((c) => c.id === v);
					const newIsAdHoc = v === "m3u" || v === "m3u8";
					const newProv = providers.find(
						(p) => p.name === newConnection?.provider,
					);
					const newResList = filterWritableResources(newProv, isTarget);
					const nextRes = newIsAdHoc
						? endpoint.resource
						: (newResList[0]?.name ?? endpoint.resource);
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
					{[...displayConnections]
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

	const resourceControl = isAdHoc ? (
		<AdHocResourceControl
			resource={endpoint.resource}
			onChange={(v) => onChange({ resource: v })}
			layout={layout}
		/>
	) : (
		<StandardResourceControl
			resource={endpoint.resource}
			availableResources={availableResources}
			isInvalidTarget={isInvalidTarget}
			onChange={(v) => onChange({ resource: v })}
			layout={layout}
		/>
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
					{resourceControl}
				</div>
				{queryInput}
			</div>
		);
	}

	return (
		<div className="flex flex-nowrap gap-2.5 items-center w-full min-w-0">
			{connectionSelect}
			{resourceControl}
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
