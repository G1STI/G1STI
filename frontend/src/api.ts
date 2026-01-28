export type Status = {
  connected: boolean;
  active_profile: string;
  updated_at: string;
  last_error: string;
  public_ip: string;
};

export type Subscription = {
  id: string;
  name: string;
  url: string;
  auto_update: string;
  last_updated: string;
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

const backend = window.backend?.App;

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
  connect: async (): Promise<void> => {
    if (!backend?.Connect) return;
    return backend.Connect();
  },
  disconnect: async (): Promise<void> => {
    if (!backend?.Disconnect) return;
    return backend.Disconnect();
  },
  refreshPublicIP: async (): Promise<string> => {
    if (!backend?.RefreshPublicIP) return "";
    return backend.RefreshPublicIP();
  },
};
