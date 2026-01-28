export {};

declare global {
  interface Window {
    go?: {
      main?: {
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
          AddSubscription: (
            name: string,
            url: string,
            autoUpdate: string
          ) => Promise<{
            id: string;
            name: string;
            url: string;
            auto_update: string;
            last_updated_ms: number;
            last_error: string;
            node_count: number;
          }>;
          UpdateSubscription: (subscriptionId: string) => Promise<void>;
          SetActiveProfile: (profileId: string) => Promise<void>;
          AddVlessURI: (tag: string, uri: string) => Promise<{
            id: string;
            name: string;
            subscription_id: string;
            node_id: string;
            tun_enabled: boolean;
            dns_mode: string;
          }>;
          Connect: () => Promise<void>;
          Disconnect: () => Promise<void>;
          RefreshPublicIP: () => Promise<string>;
        };
      };
    };
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
