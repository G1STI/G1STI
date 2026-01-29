export type Status = {
  connected: boolean;
  active_profile: string;
  updated_at_ms: number;
  last_error: string;
  public_ip: string;
};

export type Subscription = {
  id: string;
  name: string;
  url: string;
  auto_update: string;
  last_updated_ms: number;
  last_error: string;
  node_count: number;
};

export type Profile = {
  id: string;
  name: string;
  subscription_id: string;
  node_id: string;
  tun_enabled: boolean;
  dns_mode: string;
};

export type Node = {
  id: string;
  subscription_id: string;
  tag: string;
  type: string;
  server: string;
  port: number;
};

type Backend = {
  GetStatus?: () => Promise<Status>;
  ListSubscriptions?: () => Promise<Subscription[]>;
  ListProfiles?: () => Promise<Profile[]>;
  ListNodes?: (subscriptionId: string) => Promise<Node[]>;
  AddSubscription?: (
    name: string,
    url: string,
    autoUpdate: string
  ) => Promise<Subscription>;
  UpdateSubscription?: (subscriptionId: string) => Promise<void>;
  SetActiveProfile?: (profileId: string) => Promise<void>;
  AddVlessURI?: (tag: string, uri: string) => Promise<Profile>;
  Connect?: () => Promise<void>;
  Disconnect?: () => Promise<void>;
  RefreshPublicIP?: () => Promise<string>;
  GetAppLogTail?: (lines: number) => Promise<string>;
};

const backend: Backend | undefined =
  window.go?.main?.App ?? window.backend?.App;

export const api = {
  getStatus: async (): Promise<Status | null> => {
    if (!backend?.GetStatus) return null;
    return backend.GetStatus();
  },
  listSubscriptions: async (): Promise<Subscription[]> => {
    if (!backend?.ListSubscriptions) return [];
    return backend.ListSubscriptions();
  },
  listProfiles: async (): Promise<Profile[]> => {
    if (!backend?.ListProfiles) return [];
    return backend.ListProfiles();
  },
  listNodes: async (subscriptionId: string): Promise<Node[]> => {
    if (!backend?.ListNodes) return [];
    return backend.ListNodes(subscriptionId);
  },
  addSubscription: async (
    name: string,
    url: string,
    autoUpdate: string
  ): Promise<Subscription | null> => {
    if (!backend?.AddSubscription) {
      throw new Error("Backend bindings not available");
    }
    return backend.AddSubscription(name, url, autoUpdate);
  },
  updateSubscription: async (subscriptionId: string): Promise<void> => {
    if (!backend?.UpdateSubscription) {
      throw new Error("Backend bindings not available");
    }
    return backend.UpdateSubscription(subscriptionId);
  },
  setActiveProfile: async (profileId: string): Promise<void> => {
    if (!backend?.SetActiveProfile) {
      throw new Error("Backend bindings not available");
    }
    return backend.SetActiveProfile(profileId);
  },
  addVlessUri: async (tag: string, uri: string): Promise<Profile | null> => {
    if (!backend?.AddVlessURI) {
      throw new Error("Backend bindings not available");
    }
    return backend.AddVlessURI(tag, uri);
  },
  connect: async (): Promise<void> => {
    if (!backend?.Connect) {
      throw new Error("Backend bindings not available");
    }
    return backend.Connect();
  },
  disconnect: async (): Promise<void> => {
    if (!backend?.Disconnect) {
      throw new Error("Backend bindings not available");
    }
    return backend.Disconnect();
  },
  refreshPublicIP: async (): Promise<string> => {
    if (!backend?.RefreshPublicIP) {
      throw new Error("Backend bindings not available");
    }
    return backend.RefreshPublicIP();
  },
  getAppLogTail: async (lines = 200): Promise<string> => {
    if (!backend?.GetAppLogTail) {
      throw new Error("Backend bindings not available");
    }
    return backend.GetAppLogTail(lines);
  },
};
