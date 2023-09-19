import { BskyAgent } from "@atproto/api";
import { useState } from "react";
import { useAsyncEffect } from "use-async-effect";

export const useAPI = (
  username: string,
  appPassword: string,

  service: string,
  atmosfeedAPI: string,

  logout: () => void
) => {
  const [agent, setAgent] = useState<BskyAgent>();
  const [avatar, setAvatar] = useState("");
  const [loading, setLoading] = useState(true);

  useAsyncEffect(async () => {
    if (!username || !appPassword || !service) {
      setAvatar("");

      setLoading(false);

      return;
    }

    setLoading(true);

    const agent = new BskyAgent({
      service,
    });

    try {
      await agent.login({
        identifier: username,
        password: appPassword,
      });
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
    } finally {
      setLoading(false);
    }
  }, [agent]);

  const [enabled, setEnabled] = useState(false);
  const [postTTL, setPostTTL] = useState(6);
  useAsyncEffect(async () => {
    if (!avatar) {
      return;
    }

    setLoading(true);

    try {
      // TODO: Access external API here to fetch the user's existing configuration
      await new Promise((res) => setTimeout(res, 1000));

      setPostTTL(6);
      setEnabled(false);
    } catch (e) {
      console.error(e);

      logout();
    } finally {
      setLoading(false);
    }
  }, [avatar]);

  const [saveConfiguration, setSaveConfiguration] = useState(false);
  useAsyncEffect(async () => {
    if (!avatar) {
      return;
    }

    setLoading(true);

    try {
      // TODO: Access external API here to save the user's existing configuration
      await new Promise((res) => setTimeout(res, 1000));
    } catch (e) {
      console.error(e);

      logout();
    } finally {
      setLoading(false);
    }
  }, [saveConfiguration]);

  return {
    avatar,
    signedIn: avatar !== "",

    enabled,
    setEnabled,
    postTTL,
    setPostTTL,
    saveConfiguration: () => setSaveConfiguration((v) => !v),

    loading,
  };
};
