defmodule Explorer.Application do
  # See https://hexdocs.pm/elixir/Application.html
  # for more information on OTP Applications
  @moduledoc false

  use Application

  @impl true
  def start(_type, _args) do
    children = [
      ExplorerWeb.Telemetry,
      {DNSCluster, query: Application.get_env(:explorer, :dns_cluster_query) || :ignore},
      {Phoenix.PubSub, name: Explorer.PubSub},
      # Start the Finch HTTP client for sending emails
      {Finch, name: Explorer.Finch},
      # Start a worker by calling: Explorer.Worker.start_link(arg)
      # {Explorer.Worker, arg},
      # Start to serve requests, typically the last entry
      ExplorerWeb.Endpoint
    ]

    # See https://hexdocs.pm/elixir/Supervisor.html
    # for other strategies and supported options
    opts = [strategy: :one_for_one, name: Explorer.Supervisor]
    Supervisor.start_link(children, opts)
  end

  # Tell Phoenix to update the endpoint configuration
  # whenever the application is updated.
  @impl true
  def config_change(changed, _new, removed) do
    ExplorerWeb.Endpoint.config_change(changed, removed)
    :ok
  end
end

# called AlignedTask since Task is a reserved word in Elixir
defmodule AlignedTask do
  @enforce_keys [
    :verificationSystemId,
    # :proof,
    :pubInput,
    :taskCreatedBlock
  ]
  defstruct [
    :verificationSystemId,
    # :proof,
    :pubInput,
    :taskCreatedBlock
  ]
end

defmodule AlignedTaskCreatedInfo do
  @enforce_keys [:address, :block_hash, :block_number, :taskId, :transaction_hash, :aligned_task]
  defstruct [:address, :block_hash, :block_number, :taskId, :transaction_hash, :aligned_task]
end

defmodule AlignedTaskRespondedInfo do
  @enforce_keys [
    :address,
    :block_hash,
    :block_number,
    :taskId,
    :transaction_hash,
    :proofIsCorrect
  ]
  defstruct [:address, :block_hash, :block_number, :taskId, :transaction_hash, :proofIsCorrect]
end