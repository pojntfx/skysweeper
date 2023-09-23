import { IConfiguration } from "./models";

export class ConfigurationRestAPI {
  constructor(
    private apiURL: URL,
    private service: string,
    private accessJWT: string,
    private refreshJWT: string
  ) {}

  async getConfiguration(): Promise<IConfiguration> {
    const configurationURL = new URL(this.apiURL + "configuration");

    configurationURL.search = new URLSearchParams({
      service: this.service,
    }).toString();

    const response = await fetch(configurationURL.toString(), {
      headers: {
        Authorization: "Bearer " + this.accessJWT,
      },
    });

    if (response.status === 404) {
      return this.updateConfiguration({
        enabled: false,
        postTTL: 6,
      });
    }

    return response.json();
  }

  async updateConfiguration(config: IConfiguration): Promise<IConfiguration> {
    const configurationURL = new URL(this.apiURL + "configuration");

    configurationURL.search = new URLSearchParams({
      service: this.service,
    }).toString();

    return (
      await fetch(configurationURL.toString(), {
        method: "PUT",
        body: JSON.stringify(config),
        headers: {
          Authorization: "Bearer " + this.refreshJWT,
          "Content-Type": "application/json",
        },
      })
    ).json();
  }

  async deleteConfiguration() {
    const configurationURL = new URL(this.apiURL + "configuration");

    configurationURL.search = new URLSearchParams({
      service: this.service,
    }).toString();

    await fetch(configurationURL.toString(), {
      method: "DELETE",
      headers: {
        Authorization: "Bearer " + this.accessJWT,
      },
    });
  }
}
