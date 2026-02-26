use tonic::{transport::Server, Request, Response, Status};

pub mod computev1 {
    tonic::include_proto!("compute.v1");
}

use computev1::compute_service_server::{ComputeService, ComputeServiceServer};
use computev1::{SquareRequest, SquareResponse};

#[derive(Default)]
struct ComputeSvc;

#[tonic::async_trait]
impl ComputeService for ComputeSvc {
    async fn square(
        &self,
        request: Request<SquareRequest>,
    ) -> Result<Response<SquareResponse>, Status> {
        let value = request.into_inner().value;

        let square = value
            .checked_mul(value)
            .ok_or_else(|| Status::invalid_argument("value overflow when squaring"))?;

        Ok(Response::new(SquareResponse { value, square }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "0.0.0.0:50051".parse()?;

    println!("compute-engine listening on {}", addr);

    Server::builder()
        .add_service(ComputeServiceServer::new(ComputeSvc))
        .serve(addr)
        .await?;

    Ok(())
}
