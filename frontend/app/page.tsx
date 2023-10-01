"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuPortal,
  DropdownMenuSeparator,
  DropdownMenuSub,
  DropdownMenuSubContent,
  DropdownMenuSubTrigger,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { DropdownMenuLink } from "@/components/ui/dropdown-menu-link";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { PrivacyPolicy } from "@/components/ui/privacy-policy";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useToast } from "@/components/ui/use-toast";
import { useAPI } from "@/hooks/api";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Database,
  DownloadCloud,
  Laptop,
  Loader,
  LogIn,
  LogOut,
  Moon,
  MoonStar,
  Save,
  Scale,
  Sun,
  TrashIcon,
  User,
} from "lucide-react";
import { useTheme } from "next-themes";
import Image from "next/image";
import { useEffect, useState } from "react";
import { useForm, useWatch } from "react-hook-form";
import { useLocalStorage } from "usehooks-ts";
import * as z from "zod";
import logoDark from "../assets/logo-dark.png";
import logoLight from "../assets/logo-light.png";
import { Separator } from "@/components/ui/separator";

const setupFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "App password is required"),

  service: z.string().min(1, "Service is required"),
  skysweeperAPI: z.string().min(1, "SkySweeper API is required"),

  acceptedPrivacyPolicy: z.literal<boolean>(true),
});

const configurationFormSchema = z.object({
  enabled: z.boolean().optional(),
  postTTL: z.coerce
    .number()
    .int("Must be an integer")
    .positive("Most be positive"),
});

export default function Home() {
  const { setTheme } = useTheme();
  const [loginDialogOpen, setLoginDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [privacyPolicyDialogOpen, setPrivacyPolicyDialogOpen] = useState(false);

  const [username, setUsername] = useLocalStorage("skysweeper.username", "");
  const [password, setPassword] = useLocalStorage("skysweeper.password", "");

  const [service, setService] = useLocalStorage(
    "skysweeper.service",
    process.env.SKYSWEEPER_SERVICE_DEFAULT || "https://bsky.social"
  );
  const [skysweeperAPI, setSkySweeperAPI] = useLocalStorage(
    "skysweeper.skysweeperURL",
    process.env.SKYSWEEPER_API_DEFAULT || "https://api.skysweeper.p8s.lu"
  );

  const setupForm = useForm<z.infer<typeof setupFormSchema>>({
    resolver: zodResolver(setupFormSchema),
    defaultValues: {
      username,
      password,

      service,
      skysweeperAPI,

      acceptedPrivacyPolicy: false,
    },
  });

  const {
    avatar,
    did,
    signedIn,

    enabled,
    postTTL,

    saveConfiguration,
    deleteData,

    loading,
    logout,
  } = useAPI(
    username,
    password,
    service,
    skysweeperAPI,
    () => setPassword(""),
    (err, loggedOut) =>
      loggedOut
        ? toast({
            title: "You Have Been Logged Out",
            description: `Authentication with Bluesky failed and you have been logged out. The error is: "${err?.message}"`,
          })
        : toast({
            title: "An Error Occured",
            description: `An error could not be handled. The error is: "${err?.message}"`,
          })
  );

  const { setValue, ...configurationForm } = useForm<
    z.infer<typeof configurationFormSchema>
  >({
    resolver: zodResolver(configurationFormSchema),
    defaultValues: {
      enabled: false,
      postTTL: 6,
    },
  });

  const { enabled: nextEnabled, postTTL: nextPostTTL } =
    useWatch(configurationForm);

  useEffect(() => {
    setValue("enabled", enabled);
    setValue("postTTL", postTTL);
  }, [setValue, enabled, postTTL]);

  const { toast } = useToast();

  return (
    <>
      <div className="fixed w-full">
        <header className="container flex justify-between items-center py-6">
          {signedIn && (
            <Image
              src={logoDark}
              alt="SkySweeper Logo"
              className="h-10 w-auto mr-4 logo-dark"
            />
          )}

          {signedIn && (
            <Image
              src={logoLight}
              alt="SkySweeper Logo"
              className="h-10 w-auto mr-4 logo-light"
            />
          )}

          {signedIn && (
            <div className="flex content-center">
              <DropdownMenu>
                <DropdownMenuTrigger>
                  <Avatar>
                    <AvatarImage src={avatar} alt={"Avatar of " + username} />
                    <AvatarFallback>AV</AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent>
                  <DropdownMenuLabel>My Account</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuLink
                    href={`https://bsky.app/profile/${username}`}
                    target="_blank"
                  >
                    <User className="mr-2 h-4 w-4" /> Profile
                  </DropdownMenuLink>
                  <DropdownMenuItem onClick={() => logout()}>
                    <LogOut className="mr-2 h-4 w-4" /> Logout
                  </DropdownMenuItem>

                  <DropdownMenuSub>
                    <DropdownMenuSubTrigger>
                      <Database className="mr-2 h-4 w-4" />
                      <span className="mr-2">Your Data</span>
                    </DropdownMenuSubTrigger>

                    <DropdownMenuPortal>
                      <DropdownMenuSubContent>
                        <DropdownMenuItem
                          onClick={() => {
                            const data = {
                              did,
                              service,
                              enabled,
                              postTTL,
                            };

                            const blob = new Blob(
                              [JSON.stringify(data, null, 2)],
                              { type: "application/json" }
                            );
                            const url = URL.createObjectURL(blob);
                            const a = document.createElement("a");

                            a.href = url;
                            a.download = "skysweeper.json";
                            a.click();

                            URL.revokeObjectURL(url);

                            toast({
                              title: "Data Downloaded Successfully",
                              description:
                                "Your data has successfully been downloaded to your system.",
                            });
                          }}
                        >
                          <DownloadCloud className="mr-2 h-4 w-4" />
                          <span>Download your Data</span>
                        </DropdownMenuItem>
                        <DropdownMenuItem
                          onClick={() => setDeleteDialogOpen((v) => !v)}
                        >
                          <TrashIcon className="mr-2 h-4 w-4" />
                          <span>Delete your Data</span>
                        </DropdownMenuItem>
                      </DropdownMenuSubContent>
                    </DropdownMenuPortal>
                  </DropdownMenuSub>

                  <DropdownMenuLabel>Settings</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuSub>
                    <DropdownMenuSubTrigger>
                      <MoonStar className="mr-2 h-4 w-4" />
                      <span>Theme</span>
                    </DropdownMenuSubTrigger>
                    <DropdownMenuPortal>
                      <DropdownMenuSubContent>
                        <DropdownMenuItem onClick={() => setTheme("light")}>
                          <Sun className="mr-2 h-4 w-4" /> Light
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setTheme("dark")}>
                          <Moon className="mr-2 h-4 w-4" /> Dark
                        </DropdownMenuItem>
                        <DropdownMenuItem onClick={() => setTheme("system")}>
                          <Laptop className="mr-2 h-4 w-4" /> System
                        </DropdownMenuItem>
                      </DropdownMenuSubContent>
                    </DropdownMenuPortal>
                  </DropdownMenuSub>
                </DropdownMenuContent>
              </DropdownMenu>
            </div>
          )}
        </header>

        {signedIn && (
          <div className="gradient-blur">
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
            <div></div>
          </div>
        )}

        <div className="gradient-blur-bottom">
          <div></div>
          <div></div>
          <div></div>
          <div></div>
          <div></div>
          <div></div>
        </div>
      </div>

      <div className="content">
        <main className="flex-grow flex flex-col justify-center items-center gap-2 container">
          {signedIn ? (
            <>
              <Card className="max-w-full w-[500px]">
                <CardHeader>
                  <CardTitle>Configuration</CardTitle>
                </CardHeader>

                <CardContent>
                  <Form {...{ setValue, ...configurationForm }}>
                    <form
                      onSubmit={configurationForm.handleSubmit(async (v) => {
                        await saveConfiguration(
                          v.enabled ? true : false,
                          v.postTTL
                        );

                        toast({
                          title: "Configuration Saved Successfully",
                          description: v.enabled
                            ? "Your old skeets will now be deleted automatically."
                            : "Your old skeets will no longer be deleted automatically.",
                        });
                      })}
                      className="space-y-4"
                      id="configuration"
                    >
                      <FormField
                        control={configurationForm.control}
                        name="enabled"
                        render={({ field }) => {
                          const { value, onChange, ...rest } = field;

                          return (
                            <FormItem className="items-top flex space-x-2 space-y-0">
                              <FormControl>
                                <Checkbox
                                  checked={value}
                                  onCheckedChange={onChange}
                                  {...rest}
                                />
                              </FormControl>

                              <div className="grid gap-1.5 leading-none">
                                <FormLabel className="text-sm font-medium leading-none">
                                  Automatically delete old skeets
                                </FormLabel>

                                <p className="text-sm text-muted-foreground">
                                  No skeets will be deleted until you save the
                                  configuration.
                                </p>
                              </div>
                            </FormItem>
                          );
                        }}
                      />

                      {nextEnabled && (
                        <FormField
                          control={configurationForm.control}
                          name="postTTL"
                          render={({ field }) => (
                            <FormItem className="space-y-2">
                              <FormLabel>Maximum post age</FormLabel>

                              <FormDescription>
                                SkySweeper will periodically scan your skeets, and
                                if they are older than the maximum post age it
                                will delete them for you automatically.
                              </FormDescription>

                              <FormControl>
                                <div className="flex w-full items-center justify-center space-x-2 pt-2">
                                  <Input
                                    type="number"
                                    size={
                                      (nextPostTTL?.toString().length || 2) + 1
                                    }
                                    min={1}
                                    className="w-auto"
                                    {...field}
                                  />{" "}
                                  <div>
                                    month{(nextPostTTL || 0) > 1 ? "s" : ""}
                                  </div>
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      )}
                    </form>
                  </Form>
                </CardContent>

                <CardFooter>
                  <Button type="submit" form="configuration" disabled={loading}>
                    {loading ? (
                      <Loader className="mr-2 h-4 w-4 animate-spin" />
                    ) : (
                      <Save className="mr-2 h-4 w-4" />
                    )}{" "}
                    Save
                  </Button>
                </CardFooter>
              </Card>
            </>
          ) : (
            <>
              <Image
                src={logoDark}
                alt="SkySweeper Logo"
                className="h-20 w-auto logo-dark"
              />

              <Image
                src={logoLight}
                alt="SkySweeper Logo"
                className="h-20 w-auto logo-light"
              />

              <h2 className="text-2xl mt-4 my-5 text-center">
                Automatically delete your old skeets from Bluesky.
              </h2>

              <Button
                disabled={loading}
                onClick={() => setLoginDialogOpen(true)}
                className="mb-10"
              >
                {loading ? (
                  <Loader className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <LogIn className="mr-2 h-4 w-4" />
                )}{" "}
                Login with Bluesky
              </Button>
            </>
          )}
        </main>
      </div>

      <div className="fixed bottom-0 w-full overflow-x-auto">
        <footer className="flex justify-between items-center py-6 container pr-0">
          <a
            href="https://github.com/pojntfx/skysweeper"
            target="_blank"
            className="hover:underline whitespace-nowrap mr-4"
          >
            © 2023 Felicitas Pojtinger
          </a>

          <div className="flex h-5 items-center space-x-4 text-sm pr-8">
            <Button
              variant="link"
              className="p-0 h-auto font-normal"
              onClick={() => setPrivacyPolicyDialogOpen((v) => !v)}
            >
              Privacy
            </Button>

            <Separator orientation="vertical" />

            <a
              href="https://felicitas.pojtinger.com/imprint"
              target="_blank"
              className="hover:underline"
            >
              Imprint
            </a>
          </div>
        </footer>
      </div>

      <Dialog
        onOpenChange={(v) => setLoginDialogOpen(v)}
        open={loginDialogOpen}
      >
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <Image
              src={logoDark}
              alt="SkySweeper Logo"
              className="h-10 object-contain logo-dark"
            />

            <Image
              src={logoLight}
              alt="SkySweeper Logo"
              className="h-10 object-contain logo-light"
            />

            <DialogTitle className="pt-4">Login</DialogTitle>
            <DialogDescription>
              SkySweeper needs access to your Bluesky account in order to delete
              posts on your behalf.
            </DialogDescription>
          </DialogHeader>

          <Form {...setupForm}>
            <form
              onSubmit={setupForm.handleSubmit((v) => {
                setUsername(v.username);
                setPassword(v.password);

                setService(v.service);
                setSkySweeperAPI(v.skysweeperAPI);

                setLoginDialogOpen(false);
              })}
              className="space-y-4"
              id="setup"
            >
              <FormField
                control={setupForm.control}
                name="username"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Username</FormLabel>

                    <FormControl>
                      <Input type="text" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={setupForm.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>

                    <FormDescription>
                      You can use an{" "}
                      <a
                        className="underline"
                        href="https://bsky.app/settings/app-passwords"
                        target="_blank"
                      >
                        app password
                      </a>
                      . It is only stored in this browser and never uploaded to
                      our servers.
                    </FormDescription>

                    <FormControl>
                      <Input type="password" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={setupForm.control}
                name="acceptedPrivacyPolicy"
                render={({ field }) => {
                  const { value, onChange, ...rest } = field;

                  return (
                    <FormItem className="items-top flex space-x-2 space-y-0 items-center">
                      <FormControl>
                        <Checkbox
                          checked={value}
                          onCheckedChange={onChange}
                          {...rest}
                        />
                      </FormControl>

                      <div className="grid gap-1.5 leading-none">
                        <FormLabel className="text-sm font-medium leading-none">
                          I have read and agree to the{" "}
                          <Button
                            variant="link"
                            className="p-0 underline h-auto font-normal"
                            onClick={() =>
                              setPrivacyPolicyDialogOpen((v) => !v)
                            }
                          >
                            privacy policy
                          </Button>
                        </FormLabel>
                      </div>
                    </FormItem>
                  );
                }}
              />

              <Accordion type="single" collapsible>
                <AccordionItem value="item-1">
                  <AccordionTrigger>Advanced</AccordionTrigger>
                  <AccordionContent>
                    <FormField
                      control={setupForm.control}
                      name="service"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Service</FormLabel>

                          <FormDescription>
                            The Bluesky service your account is hosted on; most
                            users don&apos;t need to change this.
                          </FormDescription>

                          <FormControl>
                            <Input type="text" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={setupForm.control}
                      name="skysweeperAPI"
                      render={({ field }) => (
                        <FormItem className="mt-4">
                          <FormLabel>SkySweeper API</FormLabel>

                          <FormDescription>
                            The URL that SkySweeper&apos;s API is hosted on; most
                            users don&apos;t need to change this.
                          </FormDescription>

                          <FormControl>
                            <Input type="text" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  </AccordionContent>
                </AccordionItem>
              </Accordion>
            </form>
          </Form>

          <DialogFooter>
            <Button type="submit" form="setup">
              Next
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <AlertDialog
        onOpenChange={(v) => setDeleteDialogOpen(v)}
        open={deleteDialogOpen}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This will permanently delete your SkySweeper account and remove your
              data from our servers. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={async () => {
                await deleteData();

                toast({
                  title: "Data Deleted Successfullyy",
                  description:
                    "Your data has successfully been deleted from our servers and you have been logged out.",
                });
              }}
            >
              Delete Your Data
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <Dialog
        onOpenChange={(v) => setPrivacyPolicyDialogOpen(v)}
        open={privacyPolicyDialogOpen}
      >
        <DialogContent className="max-w-[720px] h-[720px] max-h-screen">
          <DialogHeader>
            <DialogTitle>Privacy Policy</DialogTitle>
          </DialogHeader>

          <ScrollArea className="privacy-policy">
            <PrivacyPolicy />
          </ScrollArea>
        </DialogContent>
      </Dialog>
    </>
  );
}
