import { useEffect, useMemo, useState } from "react";
import { api, Node, Profile, Status, Subscription } from "./api";

const tabs = ["Home", "Subscriptions", "Profiles", "Logs"] as const;
type Tab = (typeof tabs)[number];

const App = () => {
  const [activeTab, setActiveTab] = useState<Tab>("Home");
  const [status, setStatus] = useState<Status | null>(null);
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([]);
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [nodes, setNodes] = useState<Node[]>([]);
  const [selectedSubscription, setSelectedSubscription] = useState<string>("");
  const [search, setSearch] = useState("");
  const [subName, setSubName] = useState("");
  const [subUrl, setSubUrl] = useState("");
  const [subAutoUpdate, setSubAutoUpdate] = useState("off");
  const [vlessTag, setVlessTag] = useState("");
  const [vlessUri, setVlessUri] = useState("");
  const [formError, setFormError] = useState("");

  useEffect(() => {
    const load = async () => {
      const [statusData, subs, profs] = await Promise.all([
        api.getStatus(),
        api.listSubscriptions(),
        api.listProfiles(),
      ]);
      if (statusData) {
        setStatus(statusData);
      }
      setSubscriptions(subs);
      setProfiles(profs);
      if (subs.length > 0) {
        setSelectedSubscription(subs[0].id);
      }
    };
    void load();
  }, []);

  useEffect(() => {
    const loadNodes = async () => {
      if (!selectedSubscription) {
        setNodes([]);
        return;
      }
      const list = await api.listNodes(selectedSubscription);
      setNodes(list);
    };
    void loadNodes();
  }, [selectedSubscription]);

  const filteredNodes = useMemo(() => {
    if (!search.trim()) return nodes;
    const query = search.toLowerCase();
    return nodes.filter(
      (node) =>
        node.tag.toLowerCase().includes(query) ||
        node.server.toLowerCase().includes(query) ||
        node.type.toLowerCase().includes(query)
    );
  }, [nodes, search]);

  const onConnect = async () => {
    await api.connect();
    const nextStatus = await api.getStatus();
    if (nextStatus) setStatus(nextStatus);
  };

  const onDisconnect = async () => {
    await api.disconnect();
    const nextStatus = await api.getStatus();
    if (nextStatus) setStatus(nextStatus);
  };

  const onRefreshIP = async () => {
    await api.refreshPublicIP();
    const nextStatus = await api.getStatus();
    if (nextStatus) setStatus(nextStatus);
  };

  const onSetActiveProfile = async (profileId: string) => {
    await api.setActiveProfile(profileId);
    const nextStatus = await api.getStatus();
    if (nextStatus) setStatus(nextStatus);
  };

  const onAddSubscription = async () => {
    setFormError("");
    if (!subName.trim() || !subUrl.trim()) {
      setFormError("Name and URL are required.");
      return;
    }
    const created = await api.addSubscription(
      subName.trim(),
      subUrl.trim(),
      subAutoUpdate
    );
    if (!created) {
      setFormError("Failed to add subscription.");
      return;
    }
    setSubName("");
    setSubUrl("");
    const updated = await api.listSubscriptions();
    setSubscriptions(updated);
    setSelectedSubscription(created.id);
  };

  const onUpdateSubscription = async () => {
    if (!selectedSubscription) {
      setFormError("Select a subscription first.");
      return;
    }
    await api.updateSubscription(selectedSubscription);
    const list = await api.listNodes(selectedSubscription);
    setNodes(list);
    const updated = await api.listSubscriptions();
    setSubscriptions(updated);
  };

  const onAddVless = async () => {
    setFormError("");
    if (!vlessUri.trim()) {
      setFormError("VLESS URI is required.");
      return;
    }
    const profile = await api.addVlessUri(vlessTag.trim(), vlessUri.trim());
    if (!profile) {
      setFormError("Failed to add VLESS link.");
      return;
    }
    setVlessTag("");
    setVlessUri("");
    const [subs, profs] = await Promise.all([
      api.listSubscriptions(),
      api.listProfiles(),
    ]);
    setSubscriptions(subs);
    setProfiles(profs);
    setSelectedSubscription(profile.subscription_id);
    const list = await api.listNodes(profile.subscription_id);
    setNodes(list);
    const nextStatus = await api.getStatus();
    if (nextStatus) setStatus(nextStatus);
  };

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
              <span
                className={
                  status?.connected
                    ? "status__value status__value--ok"
                    : "status__value status__value--idle"
                }
              >
                {status?.connected ? "Connected" : "Disconnected"}
              </span>
              {status?.last_error && (
                <div className="status__error">{status.last_error}</div>
              )}
            </div>
            <div className="status__actions">
              <button className="primary" onClick={onConnect}>
                Connect
              </button>
              <button onClick={onDisconnect}>Disconnect</button>
            </div>
          </div>
          <div className="status__row">
            <label htmlFor="profile">Active profile</label>
            <select
              id="profile"
              value={status?.active_profile || ""}
              onChange={(event) => onSetActiveProfile(event.target.value)}
            >
              {profiles.length === 0 && <option>No profiles</option>}
              {profiles.map((profile) => (
                <option key={profile.id} value={profile.id}>
                  {profile.name}
                </option>
              ))}
            </select>
          </div>
          <div className="status__row">
            <label>Current IP</label>
            <div className="status__ip">
              <span className="muted">
                {status?.public_ip || "Not checked"}
              </span>
              <button onClick={onRefreshIP}>Refresh</button>
            </div>
          </div>
        </section>
      )}

      {activeTab === "Subscriptions" && (
        <section className="app__card">
          <h2>Subscriptions</h2>
          <div className="toolbar">
            <button className="primary" onClick={onAddSubscription}>
              Add subscription
            </button>
            <button onClick={onUpdateSubscription}>Update now</button>
            <input
              type="search"
              placeholder="Search nodes"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
            />
          </div>
          <div className="toolbar">
            <input
              type="text"
              placeholder="Subscription name"
              value={subName}
              onChange={(event) => setSubName(event.target.value)}
            />
            <input
              type="text"
              placeholder="Subscription URL"
              value={subUrl}
              onChange={(event) => setSubUrl(event.target.value)}
            />
            <select
              value={subAutoUpdate}
              onChange={(event) => setSubAutoUpdate(event.target.value)}
            >
              <option value="off">Auto-update: off</option>
              <option value="6h">Auto-update: 6h</option>
              <option value="12h">Auto-update: 12h</option>
              <option value="24h">Auto-update: 24h</option>
            </select>
          </div>
          <div className="toolbar">
            <input
              type="text"
              placeholder="VLESS tag (optional)"
              value={vlessTag}
              onChange={(event) => setVlessTag(event.target.value)}
            />
            <input
              type="text"
              placeholder="Paste VLESS URI"
              value={vlessUri}
              onChange={(event) => setVlessUri(event.target.value)}
            />
            <button onClick={onAddVless}>Add VLESS link</button>
          </div>
          {formError && <div className="status__error">{formError}</div>}
          <div className="toolbar">
            <label className="muted">Subscription</label>
            <select
              value={selectedSubscription}
              onChange={(event) => setSelectedSubscription(event.target.value)}
            >
              {subscriptions.length === 0 && <option>No subscriptions</option>}
              {subscriptions.map((subscription) => (
                <option key={subscription.id} value={subscription.id}>
                  {subscription.name}
                </option>
              ))}
            </select>
          </div>
          <div className="table">
            <div className="table__header">
              <span>Name</span>
              <span>Type</span>
              <span>Host</span>
              <span>Port</span>
            </div>
            {filteredNodes.length === 0 ? (
              <div className="table__row muted">No nodes loaded yet.</div>
            ) : (
              filteredNodes.map((node) => (
                <div key={node.id} className="table__row">
                  <span>{node.tag}</span>
                  <span>{node.type}</span>
                  <span>{node.server}</span>
                  <span>{node.port}</span>
                </div>
              ))
            )}
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
            {profiles.length === 0 ? (
              <div className="table__row muted">No profiles created yet.</div>
            ) : (
              profiles.map((profile) => (
                <div key={profile.id} className="table__row">
                  <span>{profile.name}</span>
                  <span>{profile.subscription_id || "-"}</span>
                  <span>{profile.node_id || "-"}</span>
                  <span>{profile.tun_enabled ? "On" : "Off"}</span>
                </div>
              ))
            )}
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
