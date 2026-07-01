import { Schema } from "effect";

/**
 * Schema for Endpoint model.
 */
export const EndpointSchema = Schema.Struct({
	source_id: Schema.String,
	resource: Schema.String,
	query: Schema.optional(Schema.String),
});

/**
 * Schema for Source model.
 */
export const SourceSchema = Schema.Struct({
	id: Schema.String,
	name: Schema.String,
	provider: Schema.String,
	config: Schema.Record({ key: Schema.String, value: Schema.String }),
});

/**
 * Schema for Step model.
 */
export const StepSchema = Schema.Struct({
	id: Schema.String,
	kind: Schema.String,
	source: EndpointSchema,
	targets: Schema.Array(EndpointSchema),
	after: Schema.optional(Schema.Array(Schema.String)),
	options: Schema.optional(
		Schema.Record({ key: Schema.String, value: Schema.Unknown }),
	),
});

/**
 * Schema for Workflow model.
 */
export const WorkflowSchema = Schema.Struct({
	id: Schema.String,
	name: Schema.String,
	steps: Schema.Array(StepSchema),
});

/**
 * Schema for TrackRow summary.
 */
export const TrackRowSchema = Schema.Struct({
	id: Schema.String,
	title: Schema.String,
	artist: Schema.String,
	bpm: Schema.String,
});

/**
 * Schema for StepDiff model.
 */
export const StepDiffSchema = Schema.Struct({
	step_id: Schema.String,
	target_name: Schema.String,
	current: Schema.NullOr(Schema.Array(TrackRowSchema)),
	added: Schema.NullOr(Schema.Array(TrackRowSchema)),
	removed: Schema.NullOr(Schema.Array(TrackRowSchema)),
	unchanged: Schema.NullOr(Schema.Array(TrackRowSchema)),
});

/**
 * Schema for ResourceInfo model.
 */
export const ResourceInfoSchema = Schema.Struct({
	name: Schema.String,
	can_write: Schema.Boolean,
	supports_query: Schema.Boolean,
});

/**
 * Schema for ProviderCapabilities model.
 */
export const ProviderCapabilitiesSchema = Schema.Struct({
	CanWrite: Schema.Boolean,
	CanManageGroups: Schema.Boolean,
	CanUpdateMetadata: Schema.Boolean,
	SupportsCues: Schema.Boolean,
	SupportsBeatgrids: Schema.Boolean,
	IsFileBased: Schema.Boolean,
});

/**
 * Schema for ProviderInfo model.
 */
export const ProviderInfoSchema = Schema.Struct({
	name: Schema.String,
	resources: Schema.Array(ResourceInfoSchema),
	capabilities: ProviderCapabilitiesSchema,
});

/**
 * Schema for GroupRow summary.
 */
export const GroupRowSchema = Schema.Struct({
	id: Schema.String,
	name: Schema.String,
	kind: Schema.String,
	parent: Schema.String,
	items: Schema.Number,
});

/**
 * Schema for QueryResult model.
 */
export const QueryResultSchema = Schema.Struct({
	kind: Schema.String,
	tracks: Schema.NullOr(Schema.Array(TrackRowSchema)),
	groups: Schema.NullOr(Schema.Array(GroupRowSchema)),
	count: Schema.Number,
});

/**
 * Schema for StepResult model.
 */
export const StepResultSchema = Schema.Struct({
	step_id: Schema.String,
	status: Schema.String,
	previews: Schema.optional(Schema.NullOr(Schema.Array(Schema.String))),
	successes: Schema.optional(Schema.NullOr(Schema.Array(Schema.String))),
	warnings: Schema.optional(Schema.NullOr(Schema.Array(Schema.String))),
	error: Schema.optional(Schema.String),
});

/**
 * Schema for WorkflowResult model.
 */
export const WorkflowResultSchema = Schema.Struct({
	workflow_id: Schema.String,
	steps: Schema.Array(StepResultSchema),
});

/**
 * Schema for UpdateConfig model.
 */
export const UpdateConfigSchema = Schema.Struct({
	last_check_at: Schema.String,
	check_interval_hour: Schema.Number,
});

/**
 * Schema for UpdateInfo model.
 */
export const UpdateInfoSchema = Schema.Struct({
	available: Schema.Boolean,
	version: Schema.optional(Schema.NullOr(Schema.String)),
	current: Schema.optional(Schema.NullOr(Schema.String)),
	release_notes: Schema.optional(Schema.NullOr(Schema.String)),
	url: Schema.optional(Schema.NullOr(Schema.String)),
});
