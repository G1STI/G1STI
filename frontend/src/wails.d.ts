export {};

declare global {
  interface Window {
    backend?: {
      App?: {
        GetStatus: () => Promise<{
          connected: boolean;
          active_profile: string;
          updated_at_ms: number;
          last_error: string;
          public_ip: string;
        }>;
        ListSubscriptions: () => Promise<
          Array<{
            id: string;
            name: string;
            url: string;
            auto_update: string;
            last_updated_ms: number;
            last_error: string;
            node_count: number;
          }>
        >;
        ListProfiles: () => Promise<
          Array<{
            id: string;
            name: string;
            subscription_id: string;
            node_id: string;
            tun_enabled: boolean;
            dns_mode: string;
          }>
        >;
        ListNodes: (subscriptionId: string) => Promise<
          Array<{
            id: string;
            subscription_id: string;
            tag: string;
            type: string;
            server: string;
            port: number;
          }>
        >;
        Connect: () => Promise<void>;
        Disconnect: () => Promise<void>;
        RefreshPublicIP: () => Promise<string>;
      };
    };
  }
}
