import { useEffect, useState } from "react";
import {
	QueryTesterControls,
	QueryTesterResults,
} from "@/components/query/QueryTester";
import type { QueryResult, Source } from "@/types";
import { ListSources, PreviewQuery } from "../../wailsjs/go/gui/App";

const asSources = (x: unknown) => (x ?? []) as Source[];
const asQueryResult = (x: unknown) => x as QueryResult;

export default function QueryTesterView() {
	const [sources, setSources] = useState<Source[]>([]);
	const [sourceID, setSourceID] = useState("");
	const [resource, setResource] = useState("tracks");
	const [query, setQuery] = useState("");
	const [result, setResult] = useState<QueryResult | null>(null);
	const [error, setError] = useState("");
	const [loading, setLoading] = useState(false);

	// biome-ignore lint/correctness/useExhaustiveDependencies: only default on mount
	useEffect(() => {
		ListSources()
			.then((s) => {
				const loaded = asSources(s);
				setSources(loaded);
				if (loaded.length > 0 && !sourceID) {
					setSourceID(loaded[0].id);
				}
			})
			.catch(() => {});
	}, []);

	// biome-ignore lint/correctness/useExhaustiveDependencies: intentionally watching value changes to clear stale results
	useEffect(() => {
		setResult(null);
		setError("");
	}, [sourceID, resource, query]);

	async function handleTest() {
		setLoading(true);
		setError("");
		setResult(null);
		try {
			setResult(asQueryResult(await PreviewQuery(sourceID, resource, query)));
		} catch (e) {
			setError(String(e));
		}
		setLoading(false);
	}

	return (
		<div className="flex flex-col h-full">
			<div className="h-14 flex items-center px-6 border-b border-border bg-[hsl(240_10%_4%)] shrink-0">
				<span className="text-sm font-semibold">Query Tester</span>
			</div>
			<div className="flex-1 overflow-hidden p-6 flex flex-col">
				<div className="flex flex-col gap-6 h-full">
					<QueryTesterControls
						sources={sources}
						sourceID={sourceID}
						resource={resource}
						query={query}
						loading={loading}
						onSourceID={setSourceID}
						onResource={setResource}
						onQuery={setQuery}
						onTest={handleTest}
					/>
					<QueryTesterResults result={result} error={error} />
				</div>
			</div>
		</div>
	);
}
