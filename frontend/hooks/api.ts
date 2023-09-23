import { ConfigurationRestAPI } from "@/api/rest";
import { BskyAgent } from "@atproto/api";
import { useCallback, useState } from "react";
import { useAsyncEffect } from "use-async-effect";

export const useAPI = (
  username: string,
  appPassword: string,

  service: string,
  aeoliusAPI: string
) => {
  const [agent, setAgent] = useState<BskyAgent>();
  const [avatar, setAvatar] = useState("");
  const [loading, setLoading] = useState(true);
  const [did, setDID] = useState("");
  const [accessJWT, setAccessJWT] = useState("");

  const logout = useCallback(() => setAccessJWT(""), []);

  useAsyncEffect(async () => {
    if (!username || !appPassword || !service) {
      setAvatar("");

      return;
    }

    setLoading(true);

    const agent = new BskyAgent({
      service,
    });

    try {
      const res = await agent.login({
        identifier: username,
        password: appPassword,
      });

      setDID(res.data.did);
      setAccessJWT(res.data.accessJwt);
    } catch (e) {
      console.error(e);

      logout();
    }

    setAgent(agent);
  }, [username, appPassword, service]);

  useAsyncEffect(async () => {
    if (!agent) {
      setAvatar("");

      return;
    }

    try {
      setAvatar(
        (
          await agent.getProfile({
            actor: username,
          })
        ).data.avatar || ""
      );
    } catch (e) {
      console.error(e);

      logout();
    }
  }, [agent]);

  const [api, setAPI] = useState<ConfigurationRestAPI>();
  useAsyncEffect(() => {
    if (!aeoliusAPI || !service || !accessJWT) {
      return;
    }

    setAPI(new ConfigurationRestAPI(new URL(aeoliusAPI), service, accessJWT));
  }, [aeoliusAPI, service, accessJWT]);

  const [enabled, setEnabled] = useState(false);
  const [postTTL, setPostTTL] = useState(6);
  useAsyncEffect(async () => {
    if (!api) {
      return;
    }

    setLoading(true);

    try {
      const res = await api.getConfiguration();

      setPostTTL(res.postTTL);
      setEnabled(res.enabled);
    } catch (e) {
      console.error(e);
    } finally {
      setLoading(false);
    }
  }, [api]);

  return {
    avatar,
    did,
    signedIn: accessJWT !== "",

    enabled,
    setEnabled,
    postTTL,
    setPostTTL,

    saveConfiguration: async () => {
      if (!accessJWT) {
        return;
      }

      setLoading(true);

      try {
        // TODO: Access external API here to save the user's existing configuration
        await new Promise((res) => setTimeout(res, 1000));
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    },
    deleteData: async () => {
      if (!accessJWT) {
        return;
      }

      setLoading(true);

      try {
        // TODO: Access external API here to delete the user configuration
        await new Promise((res) => setTimeout(res, 1000));

        logout();
      } catch (e) {
        console.error(e);
      } finally {
        setLoading(false);
      }
    },

    loading,
    logout,
  };
};
