import { useState } from "react";

const tabs = ["Home", "Subscriptions", "Profiles", "Logs"] as const;
type Tab = (typeof tabs)[number];

const App = () => {
  const [activeTab, setActiveTab] = useState<Tab>("Home");

  return (
    <div className="app">
      <header className="app__header">
        <div>
          <h1>G1STI Windows Client</h1>
          <p>Minimal client for VLESS subscriptions powered by sing-box.</p>
        </div>
        <nav className="app__tabs">
          {tabs.map((tab) => (
            <button
              key={tab}
              className={activeTab === tab ? "tab tab--active" : "tab"}
              onClick={() => setActiveTab(tab)}
            >
              {tab}
            </button>
          ))}
        </nav>
      </header>

      <section className="app__card">
        <h2>Disclaimer</h2>
        <p>
          This software is intended strictly for lawful, authorized network
          administration and corporate use. You are responsible for complying
          with local laws, policies, and provider terms. We do not endorse or
          support any misuse.
        </p>
      </section>

      {activeTab === "Home" && (
        <section className="app__card">
          <h2>Home</h2>
          <div className="status">
            <div>
              <span className="status__label">Status</span>
              <span className="status__value status__value--idle">
                Disconnected
              </span>
            </div>
            <button className="primary">Connect</button>
          </div>
          <div className="status__row">
            <label htmlFor="profile">Active profile</label>
            <select id="profile">
              <option>Default profile</option>
            </select>
          </div>
          <div className="status__row">
            <label>Current IP</label>
            <span className="muted">Not connected</span>
          </div>
        </section>
      )}

      {activeTab === "Subscriptions" && (
        <section className="app__card">
          <h2>Subscriptions</h2>
          <div className="toolbar">
            <button className="primary">Add subscription</button>
            <button>Update now</button>
            <input type="search" placeholder="Search nodes" />
          </div>
          <div className="table">
            <div className="table__header">
              <span>Name</span>
              <span>Type</span>
              <span>Host</span>
              <span>Port</span>
            </div>
            <div className="table__row muted">No nodes loaded yet.</div>
          </div>
        </section>
      )}

      {activeTab === "Profiles" && (
        <section className="app__card">
          <h2>Profiles</h2>
          <div className="toolbar">
            <button className="primary">Create profile</button>
            <button>Import JSON</button>
            <button>Export JSON</button>
          </div>
          <div className="table">
            <div className="table__header">
              <span>Name</span>
              <span>Subscription</span>
              <span>Node</span>
              <span>TUN</span>
            </div>
            <div className="table__row muted">No profiles created yet.</div>
          </div>
        </section>
      )}

      {activeTab === "Logs" && (
        <section className="app__card">
          <h2>Logs</h2>
          <div className="log">
            <p className="muted">Logs will appear here after connection.</p>
          </div>
        </section>
      )}
    </div>
  );
};

export default App;
