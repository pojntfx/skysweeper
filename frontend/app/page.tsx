"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
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
import { useAPI } from "@/hooks/api";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Laptop,
  Loader,
  LogIn,
  LogOut,
  Moon,
  MoonStar,
  Save,
  Sun,
  User,
} from "lucide-react";
import { useTheme } from "next-themes";
import Image from "next/image";
import { useEffect, useState } from "react";
import { useForm, useFormState } from "react-hook-form";
import { useLocalStorage } from "usehooks-ts";
import * as z from "zod";
import logoDark from "../assets/logo-dark.svg";
import logoLight from "../assets/logo-light.svg";

const setupFormSchema = z.object({
  username: z.string().min(1, "Username is required"),
  password: z.string().min(1, "App password is required"),

  service: z.string().min(1, "Service is required"),
  aeoliusAPI: z.string().min(1, "Aeolius API is required"),
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
  const [signInDialogOpen, setSignInDialogOpen] = useState(false);

  const [username, setUsername] = useLocalStorage("aeolius.username", "");
  const [password, setPassword] = useLocalStorage("aeolius.password", "");

  const [service, setService] = useLocalStorage(
    "aeolius.service",
    "https://bsky.social"
  );
  const [aeoliusAPI, setAeoliusAPI] = useLocalStorage(
    "aeolius.aeoliusURL",
    "https://api.aeolius.p8s.lu"
  );

  const setupForm = useForm<z.infer<typeof setupFormSchema>>({
    resolver: zodResolver(setupFormSchema),
    defaultValues: {
      username,
      password,

      service,
      aeoliusAPI,
    },
  });

  const {
    avatar,
    signedIn,

    enabled,
    setEnabled,
    postTTL,
    setPostTTL,
    saveConfiguration,

    loading,
  } = useAPI(username, password, service, aeoliusAPI, () => setPassword(""));

  const { setValue, ...configurationForm } = useForm<
    z.infer<typeof configurationFormSchema>
  >({
    resolver: zodResolver(configurationFormSchema),
    defaultValues: {
      enabled: false,
      postTTL: 6,
    },
  });

  const {
    dirtyFields: { enabled: enabledDirty },
  } = useFormState(configurationForm);

  useEffect(() => {
    setValue("enabled", enabled);
    setValue("postTTL", postTTL);
  }, [setValue, enabled, postTTL]);

  return (
    <>
      <div className="fixed w-full">
        <header className="container flex justify-between items-center py-6">
          {signedIn && (
            <Image
              src={logoDark}
              alt="Aeolius Logo"
              className="h-10 w-auto mr-4 logo-dark"
            />
          )}

          {signedIn && (
            <Image
              src={logoLight}
              alt="Aeolius Logo"
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
                  <DropdownMenuItem onClick={() => setPassword("")}>
                    <LogOut className="mr-2 h-4 w-4" /> Logout
                  </DropdownMenuItem>

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
              <Card className="w-[500px]">
                <CardHeader>
                  <CardTitle>Configuration</CardTitle>
                </CardHeader>

                <CardContent>
                  <Form {...{ setValue, ...configurationForm }}>
                    <form
                      onSubmit={configurationForm.handleSubmit((v) => {
                        setEnabled(v.enabled ? true : false);
                        setPostTTL(v.postTTL);

                        saveConfiguration();
                      })}
                      className="space-y-4"
                      id="configuration"
                    >
                      <FormField
                        disabled={loading}
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
                                  Automatically delete posts
                                </FormLabel>

                                <p className="text-sm text-muted-foreground">
                                  Enables the automatic deletion of your old
                                  skeets.
                                </p>
                              </div>
                            </FormItem>
                          );
                        }}
                      />

                      {enabledDirty && (
                        <FormField
                          disabled={loading}
                          control={configurationForm.control}
                          name="postTTL"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Maximum post age</FormLabel>

                              <FormDescription>
                                We will periodically scan your skeets, and if
                                they are older than the duration set they will
                                be automaticaly deleted.
                              </FormDescription>

                              <FormControl>
                                <Input type="number" {...field} />
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
                alt="Aeolius Logo"
                className="h-20 w-auto logo-dark"
              />

              <Image
                src={logoLight}
                alt="Aeolius Logo"
                className="h-20 w-auto logo-light"
              />

              <h2 className="text-2xl mt-4 my-5 text-center">
                Automatically delete your old skeets from Bluesky.
              </h2>

              <Button
                disabled={loading}
                onClick={() => setSignInDialogOpen(true)}
                className="mb-10"
              >
                {loading ? (
                  <Loader className="mr-2 h-4 w-4 animate-spin" />
                ) : (
                  <LogIn className="mr-2 h-4 w-4" />
                )}{" "}
                Sign in with Bluesky
              </Button>
            </>
          )}
        </main>
      </div>

      <div className="fixed bottom-0 w-full">
        <footer className="flex justify-between items-center py-6 container">
          <a
            href="https://github.com/pojntfx/aeolius"
            target="_blank"
            className="hover:underline whitespace-nowrap mr-4"
          >
            © 2023 Felicitas Pojtinger
          </a>

          <a
            href="https://felicitas.pojtinger.com/imprint"
            target="_blank"
            className="hover:underline"
          >
            Imprint
          </a>
        </footer>
      </div>

      <Dialog
        onOpenChange={(v) => setSignInDialogOpen(v)}
        open={signInDialogOpen}
      >
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <Image
              src={logoDark}
              alt="Aeolius Logo"
              className="h-10 w-auto logo-dark"
            />

            <Image
              src={logoLight}
              alt="Aeolius Logo"
              className="h-10 w-auto logo-light"
            />

            <DialogTitle className="pt-4">Sign In</DialogTitle>
            <DialogDescription>
              We need access to your Bluesky account in order to delete posts on
              your behalf.
            </DialogDescription>
          </DialogHeader>

          <Form {...setupForm}>
            <form
              onSubmit={setupForm.handleSubmit((v) => {
                setUsername(v.username);
                setPassword(v.password);

                setService(v.service);
                setAeoliusAPI(v.aeoliusAPI);

                setSignInDialogOpen(false);
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
                      name="aeoliusAPI"
                      render={({ field }) => (
                        <FormItem className="mt-4">
                          <FormLabel>Aeolius API</FormLabel>

                          <FormDescription>
                            The URL that Aeolius&apos;s API is hosted on; most
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
    </>
  );
}
