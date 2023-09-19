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

    setAvatar(
      (
        await agent.getProfile({
          actor: username,
        })
      ).data.avatar || ""
    );

    setLoading(false);
  }, [agent]);

  return {
    avatar,
    signedIn: avatar !== "",
    loading,
  };
};
