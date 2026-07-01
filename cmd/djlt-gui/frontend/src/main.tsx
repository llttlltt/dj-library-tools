import { RegistryContext } from "@effect-atom/atom-react";
import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import "./index.css";
import { registry } from "./lib/runtime";

const root = document.getElementById("app");
if (root) {
	ReactDOM.createRoot(root).render(
		<React.StrictMode>
			<RegistryContext.Provider value={registry}>
				<App />
			</RegistryContext.Provider>
		</React.StrictMode>,
	);
}
