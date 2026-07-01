// Plain-object mirrors of the Wails-generated classes.
// Use these for all React state; the bindings in wailsjs accept/return these
// shapes at runtime even though TypeScript sees them as the class types.

export interface Endpoint {
	source_id: string;
	resource: string;
	query?: string;
}

export interface Step {
	id: string;
	kind: string;
	source: Endpoint;
	targets: Endpoint[];
	after?: string[];
	options?: Record<string, unknown>;
}

export interface Workflow {
	id: string;
	name: string;
	steps: Step[];
}

export interface Source {
	id: string;
	name: string;
	provider: string;
	config: Record<string, string>;
}

export interface ProviderCapabilities {
	CanWrite: boolean;
	CanManageGroups: boolean;
	CanUpdateMetadata: boolean;
	SupportsCues: boolean;
	SupportsBeatgrids: boolean;
	IsFileBased: boolean;
}

export interface ProviderInfo {
	name: string;
	capabilities: ProviderCapabilities;
}

export interface TrackRow {
	id: string;
	title: string;
	artist: string;
	bpm: string;
}

export interface GroupRow {
	id: string;
	name: string;
	kind: string;
	parent: string;
	items: number;
}

export interface QueryResult {
	kind: "tracks" | "groups";
	tracks: TrackRow[];
	groups: GroupRow[];
	count: number;
}

export interface StepDiff {
	step_id: string;
	target_name: string;
	current: TrackRow[];
	added: TrackRow[];
	removed: TrackRow[];
	unchanged: TrackRow[];
}

export interface StepResult {
	step_id: string;
	status: string;
	previews?: string[];
	successes?: string[];
	warnings?: string[];
	error?: string;
}

export interface WorkflowResult {
	workflow_id: string;
	steps: StepResult[];
}
