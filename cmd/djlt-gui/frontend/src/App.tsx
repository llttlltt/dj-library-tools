import { FlaskConical, Music2, Workflow } from "lucide-react";
import type React from "react";
import { useState } from "react";
import { QueryTester } from "@/components/query/QueryTester";
import { cn } from "@/lib/utils";
import QueryTesterView from "./views/QueryTesterView";
import SourcesView from "./views/SourcesView";
import WorkflowsView from "./views/WorkflowsView";

type Tab = "sources" | "workflows" | "query";

const NAV: {
	id: Tab;
	label: string;
	Icon: React.FC<{ className?: string }>;
}[] = [
	{ id: "sources", label: "Sources", Icon: Music2 },
	{ id: "workflows", label: "Workflows", Icon: Workflow },
	{ id: "query", label: "Query Tester", Icon: FlaskConical },
];

export interface QueryTesterOpts {
	sourceID?: string;
	resource?: string;
	query?: string;
}

export default function App() {
	const [tab, setTab] = useState<Tab>("sources");
	const [queryTesterOpen, setQueryTesterOpen] = useState(false);
	const [queryTesterOpts, setQueryTesterOpts] = useState<QueryTesterOpts>({});

	function openQueryTester(opts?: QueryTesterOpts) {
		setQueryTesterOpts(opts ?? {});
		setQueryTesterOpen(true);
	}

	return (
		<div className="flex h-screen bg-background text-foreground overflow-hidden">
			{/* ── sidebar ───────────────────────────────────────────────── */}
			<aside className="w-52 shrink-0 flex flex-col border-r border-border bg-[hsl(240_10%_5%)]">
				<div className="h-14 flex items-center px-4 border-b border-border">
					<span className="text-sm font-semibold text-foreground">
						DJ Library Tools
					</span>
				</div>
				<nav className="flex flex-col gap-0.5 p-2 mt-1">
					{NAV.map(({ id, label, Icon }) => (
						<button
							type="button"
							key={id}
							onClick={() => setTab(id)}
							className={cn(
								"flex items-center gap-2.5 rounded-md px-3 py-2 text-sm transition-colors text-left",
								tab === id
									? "bg-accent text-accent-foreground font-medium"
									: "text-muted-foreground hover:bg-accent/50 hover:text-foreground",
							)}
						>
							<Icon className="h-4 w-4 shrink-0" />
							{label}
						</button>
					))}
				</nav>
			</aside>

			{/* ── main ──────────────────────────────────────────────────── */}
			<main className="flex-1 overflow-auto">
				{tab === "sources" && <SourcesView />}
				{tab === "workflows" && (
					<WorkflowsView onOpenQueryTester={openQueryTester} />
				)}
				{tab === "query" && <QueryTesterView />}
			</main>

			{/* ── Query Tester sheet (inline from Step editor) ─────────── */}
			<QueryTester
				open={queryTesterOpen}
				onClose={() => setQueryTesterOpen(false)}
				initialSourceID={queryTesterOpts.sourceID}
				initialResource={queryTesterOpts.resource}
				initialQuery={queryTesterOpts.query}
			/>
		</div>
	);
}
